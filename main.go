package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func logHandler(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"url":      r.RemoteAddr,
		"method":   r.Method,
		"bodySize": r.ContentLength,
		"header":   r.Header,
	}).Info("Response from func logHandler")

}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/proxy", logHandler)
	http.Handle("/", router)

	fmt.Println("Server is listening")
	http.ListenAndServe(":8080", nil)
}
