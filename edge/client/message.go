package client

import (
	"encoding/json"
	. "github.com/pigeond-io/pigeond/core"
	"github.com/gorilla/websocket"
)

type MessageType int
const (
	SUBSCRIBE MessageType = 1
	PUBLISH = 2
)

type Message struct {
	Type MessageType `json:"type"`
	Topic string `json:"topic"`
	Data string `json:"data"`
}


type MessageReader interface {
	Read (*websocket.Conn, []byte)
}


func GetMessage(messageStr []byte) *Message {
	message := &Message{}
	err := json.Unmarshal(messageStr, message)
	if err != nil {
		Error.Printf("Error in parsing message: %s Error: %s", string(messageStr), err)
	}
	return message
}
