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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//reader listens for new messages being sent to our WS endpoint
func reader(conn *websocket.Conn) error {
	for {
		// read in message
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Warning(err)
			return err
		}
		log.Println(string(p))
	}
}

//wsHandler upgrade connection to WebSocket
func wsHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	//register client
	hub.clients[conn] = true
	log.Println(hub.clients)

	err = conn.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}

	err = reader(conn)
	if err != nil {
		delete(hub.clients, conn)
	}
	log.Println(hub.clients)
}
