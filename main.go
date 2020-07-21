package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func proxyHandler(msgChan chan []byte, w http.ResponseWriter, r *http.Request) {
	//log req fields
	log.WithFields(log.Fields{
		"url":        r.Host + r.URL.Path,
		"reqestAddr": r.RemoteAddr,
		"method":     r.Method,
		"bodySize":   r.ContentLength,
		"header":     r.Header,
	}).Info("Received /proxy request")

	target := r.Header.Get("X-OXYproxy-target")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading body", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	proxyReq := &HTTPReq{
		Method:     r.Method,
		Proto:      r.Proto,
		Header:     r.Header,
		Body:       string(body),
		Host:       r.Host,
		Form:       r.Form,
		Trailer:    r.Trailer,
		RemoteAddr: r.RemoteAddr,
		Target:     target,
	}

	jsonReq, err := json.Marshal(proxyReq)
	if err != nil {
		log.Warning(err)
	}

	msgChan <- jsonReq
}

func main() {
	msgChan := make(chan []byte)
	errCh := make(chan error)
	hub := newHub()
	router := mux.NewRouter()
	router.HandleFunc("/proxy", func(w http.ResponseWriter, r *http.Request) {
		proxyHandler(msgChan, w, r)
	})
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsHandler(msgChan, errCh, hub, w, r)
	})

	log.Info("Server is listening")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
