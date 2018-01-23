// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package resp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/fzzy/radix/redis/resp"
	"io"
)

/*
	RESP - (Redis Serialization Protocol).
	http://redis.io/topics/protocol

	RESP is a compromise between the following things:

	1/ Simple to implement.
	2/ Fast to parse.
	3/ Human readable.

	RESP can serialize different data types like integers, strings, arrays. There is also a specific type for errors. Requests are sent from the client to the Redis server as arrays of strings representing the arguments of the command to execute. Redis replies with a command-specific data type.

	RESP is binary-safe and does not require processing of bulk data transferred from one process to another, because it uses prefixed-length to transfer bulk data.

	In the following example the client (C) sends the command SUBSCRIBE mytopic in order to subscribe to the topic mytopic, and the server (S) replies with OK

	C: *2\r\n
	C: $9\r\n
	C: SUBSCRIBE\r\n
	C: $7\r\n
	C: mytopic\r\n

	S: $2\r\n
	S: OK\r\n
*/

var (
	InvalidCommand = errors.New("Invalid Command")
	OkResponse     = "+OK\r\n"
	ErrorFmt       = "-Error %s\r\n"
)

type Token struct {
	Bytes []byte
}

type Command struct {
	Tokens []Token
	Err    error
}

func (cmd *Command) String() string {
	var buffer bytes.Buffer
	if cmd.Err == nil {
		buffer.WriteString("> ")
		for _, token := range cmd.Tokens {
			buffer.Write(token.Bytes)
			buffer.WriteString(" ")
		}
	} else {
		buffer.WriteString("> ERROR ")
		buffer.WriteString(cmd.Err.Error())
	}
	return buffer.String()
}

func (cmd *Command) Ok() bool {
	return cmd.Err == nil && cmd.Tokens != nil && len(cmd.Tokens) > 0
}

func (cmd *Command) Error() string {
	if cmd.Err == nil {
		return ""
	}
	return cmd.Err.Error()
}

// Make sure you call this after you have checked if command is Ok()
func (cmd *Command) Action() string {
	return string(cmd.Tokens[0].Bytes)
}

// Make sure you call this after you have checked if command is Ok()
func (cmd *Command) Args() [][]byte {
	return collectBytesFromTokens(cmd.Tokens[1:])
}

func ErrorResponse(reason string) string {
	return fmt.Sprintf(ErrorFmt, reason)
}

//Read supports Multi Commands (aka RESP Pipeline)
func Read(slice []byte) ([]*Command, bool) {
	reader := bufio.NewReader(bytes.NewReader(slice))
	ok := true
	cmds := make([]*Command, 0, 32)
	for {
		msg, err := resp.ReadMessage(reader)
		if err == io.EOF {
			break
		}
		cmd := &Command{}
		if err != nil {
			ok = false
		} else {
			ok = setCommand(cmd, msg)
		}
		cmds = append(cmds, cmd)
	}
	return cmds, ok
}

func setCommand(cmd *Command, msg *resp.Message) bool {
	ok := true
	switch msg.Type {
	case resp.Array:
		cmd.Tokens = make([]Token, 0, 32)
		tokens, err := msg.Array()
		if err == nil {
			for _, token := range tokens {
				slice, err := token.Bytes()
				if err == nil {
					cmd.Tokens = append(cmd.Tokens, Token{Bytes: slice})
				} else {
					ok = false
					cmd.Err = err
				}
			}
		} else {
			ok = false
			cmd.Err = err
		}
	case resp.SimpleStr, resp.BulkStr:
		slice, err := msg.Bytes()
		if err == nil {
			cmd.Tokens = append(cmd.Tokens, Token{Bytes: slice})
		} else {
			ok = false
			cmd.Err = err
		}
	default:
		ok = false
		cmd.Err = InvalidCommand
	}
	return ok
}

func collectBytesFromTokens(tokenSlice []Token) [][]byte {
	sliceBuf := make([][]byte, len(tokenSlice))
	for i, token := range tokenSlice {
		sliceBuf[i] = token.Bytes
	}
	return sliceBuf
}
