// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package docid

import (
	"sync"
)

/*
  Thread-safe SliceDoubleBuffer that implements DoubleBuffer
*/
type SliceDoubleBuffer struct {
	slice []DocId
	lock  sync.RWMutex //ReadWrite synchronization mutex
}

func MakeSliceDoubleBuffer() (s *SliceDoubleBuffer) {
	sDoubleBuffer := &SliceDoubleBuffer{
		slice: newSlice(),
	}
	return sDoubleBuffer
}

func (s *SliceDoubleBuffer) Add(a DocId) error {
	l := &s.lock
	l.Lock()
	s.slice = append(s.slice, a)
	l.Unlock()
	return nil
}

func (s *SliceDoubleBuffer) Slice() []DocId {
	l := &s.lock
	l.Lock()
	slice := s.slice
	s.slice = newSlice()
	l.Unlock()
	return slice
}

func newSlice() []DocId {
	return make([]DocId, 0, 256)
}
