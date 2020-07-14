package main

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func logHandler(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"url":        r.Host + r.URL.Path,
		"reqestAddr": r.RemoteAddr,
		"method":     r.Method,
		"bodySize":   r.ContentLength,
		"header":     r.Header,
	}).Info("Received /proxy request")

}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/proxy", logHandler)
	router.HandleFunc("/ws", wsHandler)

	log.Info("Server is listening")
	log.Fatal(http.ListenAndServe(":8080", router))
}
