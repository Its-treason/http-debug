package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func (s *Server) handleRedirectTo(w http.ResponseWriter, r *http.Request) {
	targetUrl := r.URL.Query().Get("url")
	if targetUrl == "" {
		errMessage := "No 'url'-parameter given!"
		w.WriteHeader(400)
		w.Write([]byte(errMessage))
		return
	}

	w.Header().Add("Location", targetUrl)
	w.WriteHeader(302)

}

func (s *Server) handleAbsoluteRedirects(w http.ResponseWriter, r *http.Request) {
	count := r.PathValue("count")
	countInt, err := strconv.ParseInt(count, 10, 32)
	if err != nil {
		w.WriteHeader(400)
		errorMessage := fmt.Sprintf("Invalid count provided! Could not be parsed as string: %s", err.Error())
		w.Write([]byte(errorMessage))
		return
	}

	countInt = countInt - 1
	if countInt <= 0 {
		w.Header().Add("Location", "/anything")
		w.WriteHeader(http.StatusFound)
		return
	}

	newUrl := *r.URL
	newUrl.Path = fmt.Sprintf("/absolute-redirect/%d", countInt)

	w.Header().Add("Location", newUrl.String())
	w.WriteHeader(http.StatusFound)
}

func (s *Server) handleRelativeRedirects(w http.ResponseWriter, r *http.Request) {
	count := r.PathValue("count")
	countInt, err := strconv.ParseInt(count, 10, 32)
	if err != nil {
		w.WriteHeader(400)
		errorMessage := fmt.Sprintf("Invalid count provided! Could not be parsed as string: %s", err.Error())
		w.Write([]byte(errorMessage))
		return
	}

	countInt = countInt - 1
	if countInt <= 0 {
		w.Header().Add("Location", "/anything")
		w.WriteHeader(http.StatusFound)
		return
	}

	w.Header().Add("Location", fmt.Sprintf("%d", countInt))
	w.WriteHeader(http.StatusFound)

}
