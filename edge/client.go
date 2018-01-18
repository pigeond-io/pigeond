// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package edge

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/pigeond-io/pigeond/common/docid"
	"github.com/pigeond-io/pigeond/common/log"
	"github.com/pigeond-io/pigeond/common/stats"
	"io"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	ConnectionClosed = 0
)

var (
	seq                int64
	emptyBuffer        = []byte{}
	ClientTickInterval = 80 * time.Millisecond
	keepAliveCounts    = (KeepAliveInterval / ClientTickInterval).Nanoseconds()
)

// WebSocketClient that encapsulates WebSocket Connection.
// For each WsClient two go routines are created.
// One goroutine for reading client requests. One goroutine for writing server updates
type WsClient struct {
	docid.StrId             // Client Id - Unique for each connection
	SessionId   docid.DocId // Each connection belongs to a unique Session
	UserId      docid.DocId // Each may connection belongs to a unique userid or is guest
	IsClosed    bool        // Is WebSocket closed
	Conn        net.Conn    // TCP based Websocket Connection
	RChan       chan int    // ClientRequestsRoutine Control Channel
	WChan       chan int    // ServerResponsesRoutine Control Channel
	once        sync.Once   // Singleton to close WebSocket once
	state       int32       // Internal State of the WsClient
}

func getNextId() string {
	return strconv.FormatInt(atomic.AddInt64(&seq, 1), 32)
}

// TODO: Generate from the Token
func getSessionId(token *jwt.Token) docid.DocId {
	return &docid.StrId{Id: "SessId"}
}

// TODO: Generate from the Token
func getUserId(token *jwt.Token) docid.DocId {
	return &docid.StrId{Id: "UserId"}
}

func InitWsClient(conn net.Conn, token *jwt.Token) {
	client := &WsClient{
		Conn:      conn,
		SessionId: getSessionId(token),
		UserId:    getUserId(token),
		IsClosed:  false,
		RChan:     make(chan int),
		WChan:     make(chan int),
		state:     0,
	}
	client.Id = getNextId()
	stats.IncrServed()
	stats.IncrLive()
	log.WithFields("edge.client", "InitWsClient").Debug("Id: ", client.SessionId, ", keepAliveCounts: ", keepAliveCounts)
	go client.wsClientRequestsProcessor()
	go client.wsServerResponsesProcessor()
}

func (client *WsClient) Close() {
	client.once.Do(func() {
		log.WithFields("edge.client", "Close").Debug("Id: ", client.SessionId)
		client.IsClosed = true
		stats.DecrLive()
		go func() {
			client.WChan <- ConnectionClosed
			client.RChan <- ConnectionClosed
		}()
	})
}

func (client *WsClient) wsClientRequestsProcessor() {
	for {
		time.Sleep(ClientTickInterval)
		select {
		case v := <-client.RChan:
			if v == ConnectionClosed {
				onConnClose(client)
				return
			}
			break
		default:
			bts, op, err := wsutil.ReadClientData(client.Conn)
			if err != nil {
				_, ok := err.(wsutil.ClosedError)
				if err == io.EOF || ok {
					client.Close()
				} else {
					log.WithFields("edge.client").Error("Id: ", client.SessionId, ", Err: ", err)
				}
				break
			}
			log.WithFields("edge.client", "request").Debug("Id: ", client.SessionId, ", Op: ", op, ", Data: ", string(bts))
			// edge.DispatchCommands(conn, cmdChannel, bts)
		}
	}
}

func (client *WsClient) wsServerResponsesProcessor() {
	var count int64
	for {
		time.Sleep(ClientTickInterval)
		count += 1
		select {
		case v := <-client.WChan:
			if v == ConnectionClosed {
				onConnClose(client)
				return
			}
			break
		default:
			// TODO
			// if my topics has updates {
			//  batch up updates to me
			// }
			if count == keepAliveCounts {
				count = 0
				log.WithFields("edge.client", "Ping").Debug("Id: ", client.SessionId)
				wsutil.WriteServerMessage(client.Conn, ws.OpPing, emptyBuffer)
			}
		}
	}
}

func onConnClose(client *WsClient) {
	state := atomic.AddInt32(&client.state, 1)
	if state == 2 {
		log.WithFields("edge.client", "Conn.Close").Debug("Id: ", client.SessionId)
		client.Conn.Close()
		close(client.WChan)
		close(client.RChan)
	}
}
