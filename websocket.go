package main

import (
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

//Hub stores connected clients
type Hub struct {
	clients map[*websocket.Conn]bool
}

func newHub() *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]bool),
	}
}

var upgrager = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrager.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	hub.clients[conn] = true
}
