package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func proxyHandler(msgChan chan []byte, w http.ResponseWriter, r *http.Request) {
	mapOfHeaders := make(map[string]string)
	mapOfTrailers := make(map[string]string)
	for hk, hv := range r.Header {
		stringHV := (strings.Join(hv, ""))
		mapOfHeaders[hk] = stringHV
	}
	for tk, tv := range r.Trailer {
		stringTV := (strings.Join(tv, ""))
		mapOfTrailers[tk] = stringTV
	}

	//log req fields
	log.WithFields(log.Fields{
		"url":        r.Host + r.URL.Path,
		"reqestAddr": r.RemoteAddr,
		"method":     r.Method,
		"bodySize":   r.ContentLength,
		"headers":    r.Header,
	}).Info("Received /proxy request")

	target := r.Header.Get("X-OXYproxy-target")
	if target == "" {
		log.Error("Target source is undefined")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 Bad Request! Destination source is not set..."))
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading body", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	proxyReq := &HTTPReq{
		Method:     r.Method,
		Proto:      r.Proto,
		Headers:    mapOfHeaders,
		Body:       string(body),
		Host:       r.Host,
		Form:       r.Form,
		Trailer:    mapOfTrailers,
		RemoteAddr: r.RemoteAddr,
		Target:     target,
	}

	jsonReq, err := json.Marshal(proxyReq)
	if err != nil {
		log.Error(err)
		return
	}

	msgChan <- jsonReq
}

func main() {
	msgChan := make(chan []byte)
	hub := newHub()
	router := mux.NewRouter()
	router.HandleFunc("/proxy", func(w http.ResponseWriter, r *http.Request) {
		proxyHandler(msgChan, w, r)
	})
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsHandler(msgChan, hub, w, r)
	})

	log.Info("Server is listening")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
