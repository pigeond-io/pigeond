package hub

import (
	"net"
	"fmt"
	"bufio"
	"log"
	"strconv"
)

type Sender struct {
	Conn net.Conn
}
var UDPSender *Sender
const port = 8002
const bufferLength = 2048

func GetSender() (*Sender, error)  {
	if UDPSender != nil {
		return UDPSender, nil
	}

	UDPSender, _ = createConnection(port)
	return UDPSender, nil
}

func (sender Sender) Close()  {
	conn := sender.Conn
	if conn == nil {
		return
	}

	conn.Close()
}

func (sender Sender) Send(message string) {
	conn := sender.Conn

	log.Print("Sending message: ", message, " connection: ", conn)
	buffer :=  make([]byte, bufferLength)
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
	conn, err := net.Dial("udp", "127.0.0.1:" + strconv.Itoa(port))
	if err != nil {
		return sender, err
	}
	sender.Conn = conn
	return sender, nil
}