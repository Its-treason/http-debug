package main

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
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

// #region Digest auth

var allowedAlgorithms = map[string]func() hash.Hash{
	"md5":     md5.New,
	"sha-256": sha256.New,
	"sha-512": sha512.New,
}

var algoToRFCName = map[string]string{
	"md5":     "MD5",
	"sha-256": "SHA-256",
	"sha-512": "SHA-512",
}

func (s *Server) handleDigestAuth(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	password := r.PathValue("password")
	algo := r.PathValue("algo")
	if algo == "" {
		algo = "md5"
	}
	algo = strings.ToLower(algo)

	realm := "testrealm@host.com"

	newHash, ok := allowedAlgorithms[algo]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Unsupported algorithm: %s. Allowed: md5, sha-256, sha-512", algo)))
		return
	}
	rfcName := algoToRFCName[algo]

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Digest ") {
		sendDigestChallenge(w, realm, rfcName)
		return
	}

	params := parseDigestHeader(strings.TrimPrefix(authHeader, "Digest "))

	required := []string{"username", "realm", "nonce", "uri", "response", "qop", "nc", "cnonce"}
	for _, key := range required {
		if params[key] == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Missing required digest parameter: %s", key)))
			return
		}
	}

	if params["username"] != username || params["realm"] != realm {
		sendDigestChallenge(w, realm, rfcName)
		return
	}

	ha1 := hashHex(newHash, fmt.Sprintf("%s:%s:%s", username, realm, password))
	ha2 := hashHex(newHash, fmt.Sprintf("%s:%s", r.Method, params["uri"]))
	expected := hashHex(newHash, fmt.Sprintf("%s:%s:%s:%s:%s:%s",
		ha1, params["nonce"], params["nc"], params["cnonce"], params["qop"], ha2,
	))

	if params["response"] != expected {
		sendDigestChallenge(w, realm, rfcName)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func sendDigestChallenge(w http.ResponseWriter, realm string, algo string) {
	nonce := generateNonce()
	opaque := generateNonce()
	header := fmt.Sprintf(
		`Digest realm="%s", nonce="%s", opaque="%s", qop="auth", algorithm=%s`,
		realm, nonce, opaque, algo,
	)
	w.Header().Set("WWW-Authenticate", header)
	w.WriteHeader(http.StatusUnauthorized)
}

func parseDigestHeader(header string) map[string]string {
	params := make(map[string]string)
	for _, part := range strings.Split(header, ",") {
		part = strings.TrimSpace(part)
		if eq := strings.IndexByte(part, '='); eq >= 0 {
			key := strings.TrimSpace(part[:eq])
			value := strings.TrimSpace(part[eq+1:])
			value = strings.Trim(value, `"`)
			params[key] = value
		}
	}
	return params
}

func hashHex(newHash func() hash.Hash, s string) string {
	h := newHash()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func generateNonce() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
