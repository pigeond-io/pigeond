package client

import (
	"github.com/gorilla/websocket"
	. "github.com/pigeond-io/pigeond/core"
)

var Clients = make(map[*websocket.Conn]bool) // connection exist or not
var TopicSubscribers = make(map[string][]*websocket.Conn) // connected clients/ topic - connection mapping


func Subscribe(conn *websocket.Conn, topicName string) error {
	Info.Print("Subscribing to topic: ", topicName)
	TopicSubscribers[topicName] = append(TopicSubscribers[topicName], conn)
	return nil
}


func Publish(topicName string, data string) error {
	if conns, ok := TopicSubscribers[topicName]; ok {
		for _, conn := range conns {
			if value, exists := Clients[conn]; (!exists) || (!value) {
				continue
			}

			Info.Println("Publishing message: ", data)
			err := conn.WriteMessage(websocket.TextMessage, []byte(data))
			if err != nil {
				Error.Println("write error: ", err)
			}
		}
	}
	return nil
}
