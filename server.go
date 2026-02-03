package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"its-treason/web-test/db"

	"github.com/charmbracelet/log"
)

// Server holds the HTTP handlers and dependencies.
type Server struct {
	db  *db.DB
	mux *http.ServeMux
}

// NewServer creates a new Server with the given database connections.
func NewServer(database *db.DB) *Server {
	s := &Server{
		db:  database,
		mux: http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.mux.HandleFunc("/anything", s.handleAnything)
	s.mux.HandleFunc("/anything/{any...}", s.handleAnything)

	s.mux.HandleFunc("DELETE /delete", s.handleAnything)
	s.mux.HandleFunc("GET /get", s.handleAnything)
	s.mux.HandleFunc("PATCH /patch", s.handleAnything)
	s.mux.HandleFunc("POST /post", s.handleAnything)
	s.mux.HandleFunc("PUT /put", s.handleAnything)
	s.mux.HandleFunc("OPTIONS /head", s.handleAnything)
	s.mux.HandleFunc("HEAD /head", s.handleAnything)

	s.mux.HandleFunc("/status/{status}", s.handleStatusCode)

	s.mux.HandleFunc("GET /user-agent", s.handleUserAgent)
	s.mux.HandleFunc("GET /ip", s.handleIp)
	s.mux.HandleFunc("GET /headers", s.handleHeader)

	s.mux.HandleFunc("/redirect-to", s.handleRedirectTo)
	s.mux.HandleFunc("GET /absolute-redirect/{count}", s.handleAbsoluteRedirects)
	s.mux.HandleFunc("GET /relative-redirect/{count}", s.handleRelativeRedirects)

	s.mux.HandleFunc("GET /_health", s.handleHealth)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	health := map[string]string{
		"status":        "ok",
		"postgres":      "ok",
		"elasticsearch": "ok",
	}

	if err := s.db.Postgres.Ping(ctx); err != nil {
		log.Error("PostgreSQL health check failed", "error", err)
		health["postgres"] = "error"
		health["status"] = "degraded"
	}

	res, err := s.db.Elasticsearch.Info()
	if err != nil {
		log.Error("Elasticsearch health check failed", "error", err)
		health["elasticsearch"] = "error"
		health["status"] = "degraded"
	} else {
		res.Body.Close()
		if res.IsError() {
			health["elasticsearch"] = "error"
			health["status"] = "degraded"
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if health["status"] != "ok" {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	json.NewEncoder(w).Encode(health)
}

func (s *Server) handleStatusCode(w http.ResponseWriter, r *http.Request) {
	wantedStatus := r.PathValue("status")
	statusInt, err := strconv.ParseInt(wantedStatus, 10, 32)
	if err != nil {
		w.WriteHeader(400)
		errorMessage := fmt.Sprintf("Invalid status code provided! Could not be parsed as string: %s", err.Error())
		w.Write([]byte(errorMessage))
		return
	}

	w.WriteHeader(int(statusInt))
}
