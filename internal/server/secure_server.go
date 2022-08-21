// Package server REST definition, instance and handlers
package server

import (
	"compress/gzip"
	"context"
	"crypto/tls"
	"encoding/json"
	stdlog "log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/iamwavecut/ct-mend/internal/config"
	"github.com/iamwavecut/ct-mend/internal/storage"
	"github.com/iamwavecut/ct-mend/resources"
	"github.com/iamwavecut/ct-mend/tools"
)

type (
	TLS struct {
		addr    string
		server  *http.Server
		timeout time.Duration
	}
	RESTHandler interface {
		WithStorageAdapter(db storage.Adapter) RESTHandler
		Select(http.ResponseWriter, *http.Request)
		Get(http.ResponseWriter, *http.Request)
		Post(http.ResponseWriter, *http.Request)
		Put(http.ResponseWriter, *http.Request)
		Delete(http.ResponseWriter, *http.Request)
	}

	ClientsHandler struct {
		db storage.Adapter
	}

	ProjectsHanlder struct {
		db storage.Adapter
	}
)

func (h *ClientsHandler) WithStorageAdapter(db storage.Adapter) RESTHandler {
	h.db = db
	return h
}

func (h *ClientsHandler) Select(w http.ResponseWriter, _ *http.Request) {
	clients, err := h.db.SelectClients()
	if !try(w, err) {
		return
	}
	err = json.NewEncoder(w).Encode(clients)
	if !try(w, err) {
		return
	}
}

func (h *ClientsHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sid, ok := vars["id"]
	if !ok {
		sid = r.URL.Query().Get("id")
	}
	ID, err := strconv.Atoi(sid)
	if !try(w, err) {
		return
	}

	client, err := h.db.GetClient(ID)
	if !try(w, err) {
		return
	}
	if client == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(client)
	if !try(w, err) {
		return
	}
}

func (h *ClientsHandler) Post(w http.ResponseWriter, r *http.Request) {
	client := storage.Client{}
	err := json.NewDecoder(r.Body).Decode(&client)
	if !try(w, err) {
		return
	}

	newClient, err := h.db.UpsertClient(&client)
	if !try(w, err) {
		return
	}
	// TODO: replace with mux named route generated path
	w.Header().Add("Location", "/clients/"+strconv.Itoa(*newClient.ID))
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(newClient)
	try(w, err)
}

//nolint:dupl
func (h *ClientsHandler) Put(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sid, ok := vars["id"]
	if !ok {
		sid = r.URL.Query().Get("id")
	}
	ID, err := strconv.Atoi(sid)
	if !try(w, err) {
		return
	}
	client := &storage.Client{}
	err = json.NewDecoder(r.Body).Decode(&client)
	if !try(w, err) {
		return
	}
	if client.ID == nil || ID != *client.ID {
		w.WriteHeader(http.StatusUnprocessableEntity)
		tools.Must(json.NewEncoder(w).Encode("validation error: path and entity ids must be equal"))
		return
	}
	updatedClient, err := h.db.UpsertClient(client)
	if !try(w, err) {
		return
	}
	// TODO: replace with mux named route generated path
	w.Header().Set("Location", "/clients/"+strconv.Itoa(*updatedClient.ID))
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(updatedClient)
	try(w, err)
}

func (h *ClientsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sid, ok := vars["id"]
	if !ok {
		sid = r.URL.Query().Get("id")
	}
	ID, err := strconv.Atoi(sid)
	if !try(w, err) {
		return
	}

	err = h.db.DeleteClient(ID)
	if !try(w, err) {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectsHanlder) WithStorageAdapter(db storage.Adapter) RESTHandler {
	h.db = db
	return h
}

func (h *ProjectsHanlder) Select(w http.ResponseWriter, _ *http.Request) {
	projects, err := h.db.SelectProjects()
	if !try(w, err) {
		return
	}
	err = json.NewEncoder(w).Encode(projects)
	if !try(w, err) {
		return
	}
}

func (h *ProjectsHanlder) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sid, ok := vars["id"]
	if !ok {
		sid = r.URL.Query().Get("id")
	}
	ID, err := strconv.Atoi(sid)
	if !try(w, err) {
		return
	}

	project, err := h.db.GetProject(ID)
	if !try(w, err) {
		return
	}
	if project == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(project)
	if !try(w, err) {
		return
	}
}

func (h *ProjectsHanlder) Post(w http.ResponseWriter, r *http.Request) {
	project := storage.Project{}
	err := json.NewDecoder(r.Body).Decode(&project)
	if !try(w, err) {
		return
	}

	newProject, err := h.db.UpsertProject(&project)
	if !try(w, err) {
		return
	}
	// TODO: replace with mux named route generated path
	w.Header().Add("Location", "/projects/"+strconv.Itoa(*newProject.ID))
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(newProject)
	try(w, err)
}

//nolint:dupl
func (h *ProjectsHanlder) Put(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sid, ok := vars["id"]
	if !ok {
		sid = r.URL.Query().Get("id")
	}
	ID, err := strconv.Atoi(sid)
	if !try(w, err) {
		return
	}
	project := &storage.Project{}
	err = json.NewDecoder(r.Body).Decode(&project)
	if !try(w, err) {
		return
	}
	if project.ID == nil || ID != *project.ID {
		w.WriteHeader(http.StatusUnprocessableEntity)
		tools.Must(json.NewEncoder(w).Encode("validation error: path and entity ids must be equal"))
		return
	}
	updatedProject, err := h.db.UpsertProject(project)
	if !try(w, err) {
		return
	}
	// TODO: replace with mux named route generated path
	w.Header().Set("Location", "/projects/"+strconv.Itoa(*updatedProject.ID))
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(updatedProject)
	try(w, err)
}

func (h *ProjectsHanlder) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sid, ok := vars["id"]
	if !ok {
		sid = r.URL.Query().Get("id")
	}
	ID, err := strconv.Atoi(sid)
	if !try(w, err) {
		return
	}

	err = h.db.DeleteProject(ID)
	if !try(w, err) {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func New(config config.TLS, db storage.Adapter, gracefulTimeout time.Duration) *TLS {
	r := mux.NewRouter().UseEncodedPath()
	r.StrictSlash(true)
	r.Use(loggingMiddleware, compressMiddleware, jsonMiddleware)

	initObjectHandler("/clients", r, (&ClientsHandler{}).WithStorageAdapter(db))
	initObjectHandler("/projects", r, (&ProjectsHanlder{}).WithStorageAdapter(db))

	cfg := &tls.Config{
		MinVersion:       tls.VersionTLS13,
		CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
		},
		GetCertificate: func(_ *tls.ClientHelloInfo) (*tls.Certificate, error) {
			crt, err := resources.FS.ReadFile("certs/server.rsa.crt")
			if err != nil {
				return nil, err
			}
			key, err := resources.FS.ReadFile("certs/server.rsa.key")
			if err != nil {
				return nil, err
			}
			certificate, err := tls.X509KeyPair(crt, key)
			return &certificate, err
		},
	}

	s := &TLS{
		addr:    config.Addr,
		timeout: gracefulTimeout,
	}

	s.server = &http.Server{
		ReadHeaderTimeout: s.timeout,
		Addr:              s.addr,
		Handler:           r,
		TLSConfig:         cfg,
		TLSNextProto:      make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		ErrorLog:          stdlog.Default(),
	}

	return s
}

func (s *TLS) Listen(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return s.server.ListenAndServeTLS("", "")
	})

	eg.Go(func() error {
		<-ctx.Done()
		timeoutCtx, cancel := context.WithTimeout(context.Background(), s.timeout)
		defer cancel()
		return s.server.Shutdown(timeoutCtx)
	})

	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

func initObjectHandler(prefix string, r *mux.Router, handler RESTHandler) {
	pr := r.PathPrefix(prefix).Subrouter()
	pr.Methods("GET").Path("/").HandlerFunc(handler.Select)
	pr.Methods("GET").Path("/{id:[0-9]+}").HandlerFunc(handler.Get)
	pr.Methods("POST").Path("/").HandlerFunc(handler.Post)
	pr.Methods("PUT").Path("/{id:[0-9]+}").HandlerFunc(handler.Put)
	pr.Methods("DELETE").Path("/{id:[0-9]+}").HandlerFunc(handler.Delete)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return handlers.LoggingHandler(log.StandardLogger().Writer(), next)
}

func compressMiddleware(next http.Handler) http.Handler {
	return handlers.CompressHandlerLevel(next, gzip.BestCompression)
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Del("Content-Type")
		w.Header().Add("Content-Type", "application/json")
		handlers.ContentTypeHandler(next, "application/json").ServeHTTP(w, r)
	})
}

func try(w http.ResponseWriter, err error) bool {
	if !tools.Try(err, true) {
		if _, ok := err.(storage.ErrNotFound); ok {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`"entity not found"`))
			return false
		}
		w.WriteHeader(http.StatusNotImplemented)
		_, _ = w.Write([]byte(`"server error: ` + err.Error() + `"`))
		return false
	}
	return true
}
