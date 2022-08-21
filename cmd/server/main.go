package main

import (
	"context"
	stdlog "log"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"golang.org/x/sync/errgroup"

	"github.com/iamwavecut/ct-mend/internal/config"
	"github.com/iamwavecut/ct-mend/internal/server"
	"github.com/iamwavecut/ct-mend/internal/storage"
	"github.com/iamwavecut/ct-mend/tools"
)

func main() {
	stdlog.SetFlags(stdlog.Lshortfile)
	stdlog.SetOutput(log.StandardLogger().Writer())
	log.SetFormatter(&prefixed.TextFormatter{
		ForceColors:     true,
		ForceFormatting: true,
	})

	// setting runtime defaults, will be overwritten by ENV vars and also by defaults defined in the struct tags
	cfg := &config.Config{
		TLS: config.TLS{
			Addr: ":8443",
		},
		Storage: config.Storage{
			Type: "sqlite",
			Addr: "./db.sqlite",
			// Type: "mongodb",
			// Addr: "mongodb://mend:mend@localhost:27017",
		},
	}
	tools.Must(env.Parse(cfg))
	log.SetLevel(cfg.AppLogLevel)

	ctx := context.WithValue(context.Background(), config.Key{}, cfg)
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill, syscall.SIGTERM)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		db, err := storage.New(ctx, &cfg.Storage)
		if !tools.Try(err) {
			return err
		}
		return server.New(cfg.TLS, cfg.GracefulTimeout).Listen(ctx, db) //nolint:wrapcheck // just no
	})
	log.Traceln("hello on", cfg.TLS.Addr, cfg.Storage.Type, cfg.Storage.Addr)
	if err := eg.Wait(); !tools.Try(err) && errors.Is(err, context.Canceled) {
		log.WithError(err).Errorln("shut down with failure")
	}
	log.Traceln("bye")
}
