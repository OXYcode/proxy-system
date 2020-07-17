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

func writer(conn *websocket.Conn, msg []byte) {
	err := conn.WriteMessage(1, msg)
	if err != nil {
		log.Println(err)
	}
}

//wsHandler upgrade connection to WebSocket
func wsHandler(jsonReqChan chan []byte, errCh chan error, hub *Hub, w http.ResponseWriter, r *http.Request) {
	greeting := []byte(`{"message": "Hi Client!"}`)
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	//register client
	hub.clients[conn] = true
	log.Println(hub.clients)

	go writer(conn, greeting)

	go reader(conn, errCh)

	go func(errCh chan error) {
		err := <-errCh
		if err != nil {
			delete(hub.clients, conn)
			log.Println(hub.clients)
		}
	}(errCh)

	go func() {
		msg := <-jsonReqChan
		writer(conn, msg)
	}()

}
