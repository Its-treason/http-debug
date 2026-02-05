package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func (s *Server) handleBasicAuth(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	password := r.PathValue("password")

	combinded := fmt.Sprintf("%s:%s", username, password)
	expectedAuth := base64.StdEncoding.EncodeToString([]byte(combinded))

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.WriteHeader(http.StatusForbidden)
		errorMessage := "No 'Authorization' header sent with request"
		w.Write([]byte(errorMessage))
		return
	}

	if !strings.HasPrefix(authHeader, "Basic ") {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Authorization header must start with 'Basic '"))
		return
	}

	authHeader = strings.TrimPrefix(authHeader, "Basic ")
	if authHeader != expectedAuth {
		w.WriteHeader(http.StatusForbidden)
		errorMessage := fmt.Sprintf("Invalid Authorization header sent! Expected '%s' got '%s'", expectedAuth, authHeader)
		w.Write([]byte(errorMessage))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleBearerAuth(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.WriteHeader(http.StatusForbidden)
		errorMessage := "No 'Authorization' header sent with request"
		w.Write([]byte(errorMessage))
		return
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Authorization header must start with 'Bearer '"))
		return
	}

	authHeader = strings.TrimPrefix(authHeader, "Bearer ")
	if authHeader != token {
		w.WriteHeader(http.StatusForbidden)
		errorMessage := fmt.Sprintf("Invalid Authorization header sent! Expected '%s' got '%s'", token, authHeader)
		w.Write([]byte(errorMessage))
		return
	}

	w.WriteHeader(http.StatusOK)
}
