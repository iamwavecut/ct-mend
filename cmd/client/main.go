package main

import (
	"bytes"
	"context"
	"crypto/tls"
	_ "embed"
	"io"
	stdlog "log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"golang.org/x/sync/errgroup"

	"github.com/iamwavecut/ct-mend/resources"
	"github.com/iamwavecut/ct-mend/tools"
)

func main() {
	stdlog.SetFlags(stdlog.Lshortfile)
	stdlog.SetOutput(log.StandardLogger().Writer())
	log.SetFormatter(&prefixed.TextFormatter{
		ForceColors:     true,
		ForceFormatting: true,
	})
	log.SetLevel(log.TraceLevel)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer cancel()
	eg, _ := errgroup.WithContext(ctx)

	eg.Go(func() error {
		host, ok := os.LookupEnv("CLIENT_HOST")
		if !ok {
			host = "127.0.0.1:8443"
		}

		raw, err := resources.FS.ReadFile("api.http")
		if !tools.Try(err) {
			return err
		}
		requests := regexp.MustCompile(`(?s)> \{%.+?%}\n`).ReplaceAllString(string(raw), "")
		requests = strings.ReplaceAll(requests, "{{host}}", host)
		requests = regexp.MustCompile(`\{\{\$randomInt}}`).ReplaceAllStringFunc(requests, func(s string) string {
			return strconv.Itoa(tools.RandInt(0, 1000))
		})
		var reqCollection []string
		for _, line := range regexp.MustCompile(`(?m)^###.*$`).Split(requests, -1) {
			reqCollection = append(reqCollection, strings.Trim(line, "\r\n\t "))
		}

		customTransport := http.DefaultTransport.(*http.Transport).Clone()
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		client := &http.Client{Transport: customTransport, Timeout: 1 * time.Second}

		reqCount := len(reqCollection) - 1
		for reqNum, reqLine := range reqCollection[1:] {
			log.Tracef("request #%d/%d", reqNum+1, reqCount)
			if reqLine == "" {
				log.Warnln("skipping")
				continue
			}
			reqParts := regexp.MustCompile(
				`(?s)(?P<method>[A-Z]+) (?P<scheme>\w+)://(?P<host>.+?)(?P<query>/.*?)\n(?P<headers>.+?)(?:\z|\n\n)(?P<body>.*\z)`,
			).FindStringSubmatch(reqLine)

			urlString := (&url.URL{
				Scheme: reqParts[2],
				Host:   reqParts[3],
				Path:   reqParts[4],
			}).String()
			if urlString == "" {
				log.Warnln("skipping")
				continue
			}

			body := &bytes.Buffer{}
			if reqParts[6] != "" {
				body = bytes.NewBuffer([]byte(reqParts[6]))
			}

			req, err := http.NewRequest(reqParts[1], urlString, body)
			if !tools.Try(err) {
				return err
			}
			for _, headerLine := range strings.Split(reqParts[5], "\n") {
				header := strings.Split(headerLine, ": ")

				req.Header.Set(header[0], header[1])
			}
			log.Debugln(reqParts[1], urlString)
			resp, err := client.Do(req)
			if !tools.Try(err) {
				return err
			}
			if resp.StatusCode >= 200 && resp.StatusCode < 400 {
				log.Infoln(resp.Status, resp.StatusCode)
			} else {
				log.Warnln(resp.Status, resp.StatusCode)
			}

			b, err := io.ReadAll(resp.Body)
			if !tools.Try(err) {
				log.Errorln("error while request", reqLine)
				return err
			}
			log.Info(strings.Trim(spew.Sdump(string(b)), "\r\n\t "))
		}

		return nil
	})

	if err := eg.Wait(); !tools.Try(err) && err != context.Canceled {
		log.WithError(err).Fatalln("context was canceled with error")
	}
}
