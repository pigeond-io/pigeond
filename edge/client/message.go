// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package client

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	. "github.com/pigeond-io/pigeond/common"
)

type MessageType int

const (
	SUBSCRIBE MessageType = 1
	PUBLISH               = 2
)

type Message struct {
	Type  MessageType `json:"type"`
	Topic string      `json:"topic"`
	Data  string      `json:"data"`
}

type MessageReader interface {
	Read(*websocket.Conn, []byte)
}

func GetMessage(messageStr []byte) *Message {
	message := &Message{}
	err := json.Unmarshal(messageStr, message)
	if err != nil {
		Error.Printf("Error in parsing message: %s Error: %s", string(messageStr), err)
	}
	return message
}
