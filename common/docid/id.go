// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package docid

import (
	"strconv"
)

var (
	emptyString = ""
	eof         = "EOF Error"
)

type Nil struct {
}

type IntId struct {
	Id    int64
	docid string
}

type ByteId struct {
	Id    []byte
	docid string
}

type StrId struct {
	Id string
}

func (s *Nil) DocId() string {
	return emptyString
}

func (s *Nil) Error() string {
	return eof
}

func (i *IntId) DocId() string {
	if i.docid == "" {
		i.docid = strconv.FormatInt(i.Id, 10)
	}
	return i.docid
}

func (i *IntId) ClearDocId() {
	i.docid = ""
}

func (i *ByteId) DocId() string {
	if i.docid == "" {
		i.docid = string(i.Id)
	}
	return i.docid
}

func (i *ByteId) ClearDocId() {
	i.docid = ""
}

func (i *StrId) DocId() string {
	return i.Id
}
