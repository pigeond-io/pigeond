package message

import (
	. "github.com/pigeond-io/pigeond/core"
	"github.com/gorilla/websocket"
	"github.com/pigeond-io/pigeond/edge/hub"
	. "github.com/pigeond-io/pigeond/edge/client"
)

type DefaultMessageReader struct {

}

var hubSender *hub.Sender

func (reader DefaultMessageReader) Read(conn *websocket.Conn, messageBytes []byte) {
	Info.Println("Received message: ", string(messageBytes))

	message := GetMessage(messageBytes)
	Info.Println("Received message from topic: ", message.Topic, " of type: ", string(messageBytes), "  contains data: ", message.Data)

	switch message.Type {
	case SUBSCRIBE:
		Subscribe(conn, message.Topic)
		break
	case PUBLISH:
		hubSender, _ = hub.GetSender()
		hubSender.Send(string(messageBytes))
		break
	}
}

