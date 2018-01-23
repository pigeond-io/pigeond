// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package edge

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gobwas/ws"
	"github.com/pigeond-io/pigeond/common/docid"
	"github.com/pigeond-io/pigeond/common/log"
	"github.com/pigeond-io/pigeond/common/stats"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

var (
	KeepAliveInterval         = 1 * time.Minute
	allowAnonymousConnections = true
	jwtSecretKey              = []byte("PigeondJWTSecretKey")
)

const (
	SessionIdx int = iota
	UserIdx
	TopicIdx
)

type WsServer struct {
	indexMap docid.ImmutableIndexMap
	listener net.Listener
}

// Zero-copy Upgrade Websocket Server which allows both anonymous and jwt based authorized connections.
// There is no method in the JavaScript WebSockets API for specifying additional headers for the client/browser to send
// WsServer uses request uri path as the JwtToken
func InitWsServer(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.WithFields("edge.server").Fatal(err)
	}
	server := &WsServer{
		indexMap: docid.MakeImmutableIndexMap(SessionIdx, UserIdx, TopicIdx),
		listener: listener,
	}
	server.acceptWsClients()
}

// Public interface for clients to perform action on Server Index
func (server *WsServer) OnIndex(indexActionCallback func(docid.ImmutableIndexMap)) {
	indexActionCallback(server.indexMap)
}

// Server run loop that accepts new client connections
func (server *WsServer) acceptWsClients() {
	listener := server.listener
	for {
		conn, err := listener.Accept()
		if err != nil {
			stats.IncrFailed()
			log.WithFields("edge.server").Error(err)
		} else {
			go server.initWsClient(conn)
		}
	}
}

// Initiating a Websocket Connection
// This method enables tcp keep alive, upgrades the connection to websocket, parses and validates the jwt token if provided and initiate the wsclient
func (server *WsServer) initWsClient(conn net.Conn) {
	tcp, ok := conn.(*net.TCPConn)
	if ok {
		tcp.SetKeepAlive(true)
		tcp.SetKeepAlivePeriod(KeepAliveInterval)
	} else {
		log.WithFields("edge.server").Error("KeepAliveFailed")
	}
	var token string
	wsUpgrader := ws.Upgrader{
		OnRequest:       onWsUpgradeRequest(&token),
		OnBeforeUpgrade: beforeWsUpgrade,
	}
	_, err := wsUpgrader.Upgrade(conn)
	if err != nil {
		terminateConnection(conn, err)
		return
	}
	if token == "" {
		if !allowAnonymousConnections {
			terminateConnection(conn, "Anonymous Connections Not Allowed")
		} else {
			InitWsClient(server, conn, nil)
		}
	} else {
		jToken, err := parseToken(token)
		if err != nil {
			terminateConnection(conn, err)
			return
		}
		if !jToken.Valid {
			terminateConnection(conn, "Invalid Token")
			return
		}
		InitWsClient(server, conn, jToken)
	}
}

// TODO: Add the Host Check
func isHostOk(host string) bool {
	return true
}

// Terminates the connection on Error
func terminateConnection(conn net.Conn, err interface{}) {
	stats.IncrFailed()
	log.WithFields("edge.server").Debug(err)
	conn.Close()
	return
}

// Parses the JWT token
// TODO: Define & Validate Manadatory JWT Claims Fields
func parseToken(token string) (*jwt.Token, error) {
	log.WithFields("edge.server").Debug("Token: %s", token)
	jToken, err := jwt.Parse(token,
		func(jToken *jwt.Token) (interface{}, error) {
			if _, ok := jToken.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", jToken.Header["alg"])
			}
			return jwtSecretKey, nil
		})
	return jToken, err
}

// Before WebSocket Uprade OnRequest Callback
// Here we check the Host and Request Uri. If everything is okay we store the JWT token from the request uri.
func onWsUpgradeRequest(token *string) func([]byte, []byte) (error, int) {
	return func(host, uri []byte) (err error, code int) {
		if !isHostOk(string(host)) {
			return fmt.Errorf("Bad Request"), 403
		}
		urlObj, err := url.Parse(string(uri))
		if err == nil {
			*token = urlObj.Path[1:] //remove the forward slash
			return
		} else {
			return fmt.Errorf("Bad Request"), 403
		}
	}
}

// Before WebSocket Uprade OnUpgrade Callback
// We modify the response headers to add X-Server tag
func beforeWsUpgrade() (headerWriter func(io.Writer), err error, code int) {
	header := http.Header{
		"X-Server": []string{"pigeond-ws"},
	}
	return ws.HeaderWriter(header), nil, 0
}
