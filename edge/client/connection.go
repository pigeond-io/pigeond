// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package client

import (
	"github.com/gorilla/websocket"
	"github.com/pigeond-io/pigeond/common/log"
)

var Clients = make(map[*websocket.Conn]bool)              // connection exist or not
var TopicSubscribers = make(map[string][]*websocket.Conn) // connected clients/ topic - connection mapping

func Subscribe(conn *websocket.Conn, topicName string) error {
	log.Info("Subscribing to topic: ", topicName)
	TopicSubscribers[topicName] = append(TopicSubscribers[topicName], conn)
	return nil
}

func Publish(topicName string, data string) error {
	if conns, ok := TopicSubscribers[topicName]; ok {
		for _, conn := range conns {
			if value, exists := Clients[conn]; (!exists) || (!value) {
				continue
			}

			log.Info("Publishing message: ", data)
			err := conn.WriteMessage(websocket.TextMessage, []byte(data))
			if err != nil {
				log.Error("write error: ", err)
			}
		}
	}
	return nil
}
