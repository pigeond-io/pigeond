// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package docid

import(
	"strconv"
)

var (
	emptyString = ""
	eof = "EOF Error"
)

type Nil struct {
}

type IntId struct {
	Id int64
}

type ByteId struct {
	Id []byte
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
	return strconv.FormatInt(i.Id, 36)
}

func (i *ByteId) DocId() string {
	return string(i.Id)
}

func (i *StrId) DocId() string {
	return i.Id
}
