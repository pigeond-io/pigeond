// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package hub

import (
	"net"
	"fmt"
	. "github.com/pigeond-io/pigeond/common"
	"github.com/pigeond-io/pigeond/edge/client"
)


func Listen(port int, buffer int) {
	conn, err := connect(port)
	if err != nil {
		Error.Println("Error in starting connection: ", err)
		return
	}

	messageBytes := make([]byte, buffer) // buffer default size 2048

	for {
		n,remoteaddr,err := conn.ReadFromUDP(messageBytes)
		fmt.Printf("Read a message from %v %s \n", remoteaddr, string(messageBytes[:n]))
		if err !=  nil {
			fmt.Printf("Some error  %v", err)
			continue
		}

		message := client.GetMessage(messageBytes[:n])
		client.Publish(message.Topic, message.Data)
	}
}

func connect(port int) (*net.UDPConn, error)  {
	addr := net.UDPAddr{
		Port: port,
		IP: net.ParseIP("127.0.0.1"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}