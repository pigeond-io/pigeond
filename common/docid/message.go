// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package docid

import (
	"time"
)

type Message struct {
	StrId
	Source    DocId
	Content   []byte
	Timestamp int64
}

func MakeMessage(source DocId, content []byte) *Message {
	msg := &Message{Source: source, Content: content, Timestamp: time.Now().Unix()}
	msg.Id = MD5([]byte(source.DocId()), content)
	return msg
}

func (m *Message) Body() []byte {
	return m.Content
}
