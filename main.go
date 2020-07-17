package main

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"url":        r.Host + r.URL.Path,
		"reqestAddr": r.RemoteAddr,
		"method":     r.Method,
		"bodySize":   r.ContentLength,
		"header":     r.Header,
	}).Info("Received /proxy request")

}

func main() {
	errCh := make(chan error)
	hub := newHub()
	router := mux.NewRouter()
	router.HandleFunc("/proxy", proxyHandler)
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsHandler(errCh, hub, w, r)
	})

	log.Info("Server is listening")
	log.Fatal(http.ListenAndServe(":8080", router))
}
