package client

import (
	"github.com/gorilla/websocket"
	"net/http"
	. "github.com/pigeond-io/pigeond/core"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Handler(w http.ResponseWriter, r *http.Request, reader MessageReader) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		Error.Print("upgrade Error:", err)
		return
	}

	Clients[conn] = true
	defer func() {
		Clients[conn] = false
		conn.Close()
	}()

	for {
		messageType, messageStr, err := conn.ReadMessage()
		if err != nil {
			Error.Println("read Error:", err)
			break
		}

		if messageType != websocket.TextMessage {
			continue
		}

		reader.Read(conn, messageStr)
	}

}
