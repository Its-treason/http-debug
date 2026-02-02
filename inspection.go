package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/charmbracelet/log"
)

type AnyThingResponse struct {
	Method        string      `json:"method"`
	Args          url.Values  `json:"args"`
	Data          string      `json:"data"`
	Json          any         `json:"json"`
	Form          url.Values  `json:"form"`
	Headers       http.Header `json:"headers"`
	RemoteAddress string      `json:"remoteAddress"`
	Url           string      `json:"url"`
	Host          string      `json:"host"`
}

func (s *Server) handleAnything(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Warn("Could not read body", err)
	}
	dataString := string(data)

	var parsedJson any
	decoder := json.NewDecoder(strings.NewReader(dataString))
	decoder.UseNumber()
	err = decoder.Decode(&parsedJson)
	if err != nil {
		parsedJson = nil
	}

	// Replace body so ParseMultipartForm can read it
	r.Body = io.NopCloser(bytes.NewReader(data))
	// Ignore
	r.ParseMultipartForm(1e+9)

	jsonEncoded, err := json.Marshal(AnyThingResponse{
		Method:        r.Method,
		Args:          r.URL.Query(),
		Data:          dataString,
		Json:          parsedJson,
		Form:          r.PostForm,
		Headers:       r.Header,
		RemoteAddress: r.RemoteAddr,
		Url:           r.RequestURI,
		Host:          r.Host,
	})
	if err != nil {
		errMessage := fmt.Sprintf("Failed to encode JSON: %s", err)
		w.Write([]byte(errMessage))
		w.WriteHeader(500)
		return
	}
	w.Header().Add("content-type", "application/json")
	w.Write(jsonEncoded)
}

func (s *Server) handleUserAgent(w http.ResponseWriter, r *http.Request) {
	userAgent := r.Header.Get("user-agent")
	w.Write([]byte(userAgent))
}

func (s *Server) handleHeader(w http.ResponseWriter, r *http.Request) {
	headers, err := json.Marshal(r.Header)
	if err != nil {
		errMessage := fmt.Sprintf("Failed to encode JSON: %s", err)
		w.Write([]byte(errMessage))
		w.WriteHeader(500)
		return
	}

	w.Write(headers)
}

func (s *Server) handleIp(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr
	w.Write([]byte(ip))
}
