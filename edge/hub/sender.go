// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package hub

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
)

type Sender struct {
	Conn net.Conn
}

var UDPSender *Sender

const port = 8002
const bufferLength = 2048

func GetSender() (*Sender, error) {
	if UDPSender != nil {
		return UDPSender, nil
	}

	UDPSender, _ = createConnection(port)
	return UDPSender, nil
}

func (sender Sender) Close() {
	conn := sender.Conn
	if conn == nil {
		return
	}

	conn.Close()
}

func (sender Sender) Send(message string) {
	conn := sender.Conn

	log.Print("Sending message: ", message, " connection: ", conn)
	buffer := make([]byte, bufferLength)
	fmt.Fprintf(conn, message)

	_, err := bufio.NewReader(conn).Read(buffer)
	if err == nil {
		fmt.Printf("%s\n", buffer)
	} else {
		fmt.Printf("Some error %v\n", err)
	}
}

func createConnection(port int) (*Sender, error) {
	sender := &Sender{}
	conn, err := net.Dial("udp", "127.0.0.1:"+strconv.Itoa(port))
	if err != nil {
		return sender, err
	}
	sender.Conn = conn
	return sender, nil
}
