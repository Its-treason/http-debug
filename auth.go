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
	encoded := base64.StdEncoding.EncodeToString([]byte(combinded))

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
	if authHeader != encoded {
		w.WriteHeader(http.StatusForbidden)
		errorMessage := fmt.Sprintf("Invalid Authorization header sent! Expected '%s' got '%s'", encoded, authHeader)
		w.Write([]byte(errorMessage))
		return
	}

	w.WriteHeader(http.StatusOK)
}
