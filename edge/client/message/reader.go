// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package message

import (
	"github.com/gorilla/websocket"
	. "github.com/pigeond-io/pigeond/common"
	. "github.com/pigeond-io/pigeond/edge/client"
	"github.com/pigeond-io/pigeond/edge/hub"
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
