package main

import (
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

//Hub stores connected clients
type Hub struct {
	clients     map[*websocket.Conn]bool
	broadcaster chan []byte //broadcaster chan to share msgs between connected clients
}

func newHub() *Hub {
	return &Hub{
		clients:     make(map[*websocket.Conn]bool),
		broadcaster: make(chan []byte),
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//reader listens for new messages being sent to our WS endpoint
func reader(conn *websocket.Conn, errCh chan error) {
	for {
		// read in message
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Warning(err)
			errCh <- err
			return
		}
		log.Println(string(p))
	}
}

func writer(hub *Hub, msgChan chan []byte, broadcaster chan []byte) {
	for {
		select {
		case msg := <-msgChan:
			//sending message to all clients from Hub
			for client := range hub.clients {
				if err := client.WriteMessage(1, msg); err != nil {
					log.Println(err)
				}
			}
			break
		case greeting := <-broadcaster:
			// sending message to all clients from Hub
			for client := range hub.clients {
				if err := client.WriteMessage(1, greeting); err != nil {
					log.Println(err)
				}
			}
			break
		}
	}
}

//wsHandler upgrade connection to WebSocket
func wsHandler(msgChan chan []byte, errCh chan error, hub *Hub, w http.ResponseWriter, r *http.Request) {
	greeting := []byte(`{"message": "Hello to new Client!"}`)
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	//register client
	hub.clients[conn] = true
	log.Println(hub.clients)

	go writer(hub, msgChan, hub.broadcaster)

	go reader(conn, errCh)

	go func(errCh chan error) {
		err := <-errCh
		if err != nil {
			delete(hub.clients, conn)
			log.Println(hub.clients)
		}
	}(errCh)
	hub.broadcaster <- greeting
}
