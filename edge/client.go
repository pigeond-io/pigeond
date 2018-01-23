// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package edge

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/pigeond-io/pigeond/common/commands"
	"github.com/pigeond-io/pigeond/common/docid"
	"github.com/pigeond-io/pigeond/common/log"
	"github.com/pigeond-io/pigeond/common/resp"
	"github.com/pigeond-io/pigeond/common/stats"
	"github.com/pigeond-io/pigeond/edge/actions"
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
	cmdRegistry commands.Registry
	once        sync.Once // Singleton to close WebSocket once
	state       int32     // Internal State of the WsClient
	server      *WsServer
}

func InitWsClient(server *WsServer, conn net.Conn, token *jwt.Token) {
	var claims jwt.MapClaims
	claims = nil
	if token != nil {
		claims, _ = token.Claims.(jwt.MapClaims)
	}
	connId := getNextId()
	client := &WsClient{
		Conn:        conn,
		SessionId:   getSessionId(claims, connId),
		UserId:      getUserId(claims),
		IsClosed:    false,
		RChan:       make(chan int),
		WChan:       make(chan int),
		cmdRegistry: commands.MakeRegistry(),
		state:       0,
		server:      server,
	}
	client.Id = connId
	client.registerCommands()
	client.registerSession()
	client.registerUser()
	stats.IncrServed()
	stats.IncrLive()
	log.WithFields("edge.client", "InitWsClient").Debug(client.String(), ", keepAliveCounts: ", keepAliveCounts)
	go client.wsClientRequestsProcessor()
	go client.wsServerResponsesProcessor()
}

func (client *WsClient) String() string {
	return fmt.Sprintf("WsClient #%s", client.DocId())
}

func (client *WsClient) IsGuestSession() bool {
	return docid.Equals(client, client.SessionId)
}

func (client *WsClient) Subscribe(topic string) bool {
	log.WithFields("edge.client", "Subscribe", topic).Debug(client.String())
	//TODO
	return true
}

func (client *WsClient) Unsubscribe(topic string) bool {
	log.WithFields("edge.client", "Unsubscribe", topic).Debug(client.String())
	//TODO
	return true
}

func (client *WsClient) Close() {
	client.once.Do(func() {
		log.WithFields("edge.client", "Close").Debug(client.String())
		client.IsClosed = true
		stats.DecrLive()
		go func() {
			client.WChan <- ConnectionClosed
			client.RChan <- ConnectionClosed
		}()
	})
}

func (client *WsClient) registerCommands() {
	registry := client.cmdRegistry
	registry.Write("SUBSCRIBE", actions.OnSubscribe(client))
	registry.Write("UNSUBSCRIBE", actions.OnSubscribe(client))
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
			bts, _, err := wsutil.ReadClientData(client.Conn)
			if err != nil {
				_, ok := err.(wsutil.ClosedError)
				if err == io.EOF || ok {
					client.Close()
				} else {
					log.WithFields("edge.client").Error(client.String(), ", Err: ", err)
				}
				break
			}
			go client.executeClientRequest(bts)
		}
	}
}

// RESP Based Command Executor
func (client *WsClient) executeClientRequest(commandBytes []byte) {
	log.WithFields("edge.clientRequest").Debug(client.String(), string(commandBytes))
	cmds, ok := resp.Read(commandBytes)
	if ok {
		for _, cmd := range cmds {
			response := resp.OkResponse
			if cmd.Ok() {
				executor := commands.MakeExecutor(cmd)
				result := executor.Execute(client.cmdRegistry)
				if result != nil {
					response = resp.ErrorResponse(result.Error())
				}
			} else {
				response = resp.ErrorResponse(cmd.Error())
			}
			wsutil.WriteServerMessage(client.Conn, ws.OpText, []byte(response))
		}
	} else {
		wsutil.WriteServerMessage(client.Conn, ws.OpText, []byte(resp.ErrorResponse("Parsing Failed")))
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
				log.WithFields("edge.client", "Ping").Debug(client.String())
				wsutil.WriteServerMessage(client.Conn, ws.OpPing, emptyBuffer)
			}
		}
	}
}

func (client *WsClient) incrUserCount() {
	//TODO Track User Counts
	//If user is in dirty list. remove user from it.
}

func (client *WsClient) decrUserCount() {
	//TODO Track User Counts
	//If user count is zero then add the user to dirty list
	//Dirty sessions will expire after timeout. When they expire user is unsubscribed from all the topics subscribed.
}

func (client *WsClient) incrSessionCount() {
	//TODO Track Session Counts
	//If session is in dirty list. remove session from it.
}

func (client *WsClient) decrSessionCount() {
	//TODO Track Session Counts
	//If session count is zero then add the session to dirty list
	//Dirty sessions will expire after timeout. When they expire session is unsubscribed from all the topics subscribed.
}

// Helper function that wraps OnIndex call on the server
func (client *WsClient) onIndex(indexActionCallback func(docid.ImmutableIndexMap)) {
	server := client.server
	if server != nil {
		server.OnIndex(indexActionCallback)
	}
}

// Adds Session to SessionIdx and Does Session Management
func (client *WsClient) registerSession() {
	client.onIndex(func(index docid.ImmutableIndexMap) {
		index.Add(SessionIdx, func(idx docid.AddIndexEntryWriter) error {
			return idx.Add(client.SessionId, client)
		})
	})
	if !client.IsGuestSession() {
		client.incrSessionCount()
	}
}

// Adds User to UserIdx and Does User Management
func (client *WsClient) registerUser() {
	if docid.IsNil(client.UserId) {
		return
	}
	client.onIndex(func(index docid.ImmutableIndexMap) {
		index.Add(UserIdx, func(idx docid.AddIndexEntryWriter) error {
			return idx.Add(client.UserId, client)
		})
	})
	if !client.IsGuestSession() {
		client.incrUserCount()
	}
}

// Removes Session from SessionIdx and Does Session Management
func (client *WsClient) deregisterSession() {
	client.onIndex(func(index docid.ImmutableIndexMap) {
		index.Remove(SessionIdx, func(idx docid.RemoveIndexEntryWriter) error {
			return idx.Remove(client.SessionId, client)
		})
		if client.IsGuestSession() {
			index.RemoveValue(TopicIdx, client)
		} else {
			client.decrSessionCount()
		}
	})
}

// Removes User from UserIdx and Does User Management
func (client *WsClient) deregisterUser() {
	if docid.IsNil(client.UserId) {
		return
	}
	client.onIndex(func(index docid.ImmutableIndexMap) {
		index.Remove(UserIdx, func(idx docid.RemoveIndexEntryWriter) error {
			return idx.Remove(client.UserId, client)
		})
	})
	if !client.IsGuestSession() {
		client.decrUserCount()
	}
}

func getNextId() string {
	return strconv.FormatInt(atomic.AddInt64(&seq, 1), 32)
}

func getSessionId(claims jwt.MapClaims, connId string) docid.DocId {
	var docId docid.DocId
	docId = &docid.StrId{Id: connId}
	if claims != nil {
		sid, ok := claims["sid"].(string)
		if ok {
			docId = &docid.StrId{Id: sid}
		}
	}
	return docId
}

func getUserId(claims jwt.MapClaims) docid.DocId {
	var docId docid.DocId
	docId = &docid.Nil{}
	if claims != nil {
		uid, ok := claims["uid"].(string)
		if ok {
			docId = &docid.StrId{Id: uid}
		}
	}
	return docId
}

// Deinit method invoked when client connection is terminated.
func onConnClose(client *WsClient) {
	state := atomic.AddInt32(&client.state, 1)
	if state == 2 {
		log.WithFields("edge.client", "Conn.Close").Debug(client.String())
		client.Conn.Close()
		close(client.WChan)
		close(client.RChan)
		client.deregisterSession()
		client.deregisterUser()
		client.cmdRegistry.Close()
		client.cmdRegistry = nil
		client.server = nil
	}
}
