// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package docid

import (
	"errors"
	"sync"
)

/*
  Thread-safe HashMap that implements Map Interface
*/
type HashMap struct {
	index map[string]DocId //index map to map Id to DocId
	lock  sync.RWMutex          //ReadWrite synchronization mutex
}

// Constructor to create hashmap.
// Optionally you can pass in the initial set of docIds to populate the map
func MakeHashMap(docIds ...DocId) *HashMap {
	hashmap := &HashMap{
		index: make(map[string]DocId),
	}
	for _, docId := range docIds {
		hashmap.index[docId.Id()] = docId
	}
	return hashmap
}

func (s *HashMap) Add(a DocId) {
	l := &s.lock
	l.RLock()
	id := a.Id()
	_, ok := s.index[id]
	if !ok {
		l.RUnlock()
		l.Lock()
		s.index[id] = a
		l.Unlock()
	} else {
		l.RUnlock()
	}
}

func (s *HashMap) Remove(id string) {
	l := &s.lock
	l.RLock()
	_, ok := s.index[id]
	if ok {
		l.RUnlock()
		l.Lock()
		delete(s.index, id)
		l.Unlock()
	} else {
		l.RUnlock()
	}
}

func (s *HashMap) Get(id string) (DocId, error) {
	l := &s.lock
	l.RLock()
	docId, ok := s.index[id]
	l.RUnlock()
	if !ok {
		return EOF, errors.New("doc_id not found")
	}
	return docId, nil
}

func (s *HashMap) Contains(id string) bool {
	l := &s.lock
	l.RLock()
	_, ok := s.index[id]
	l.RUnlock()
	return ok
}

func (s *HashMap) Clear() {
	l := &s.lock
	l.Lock()
	s.index = make(map[string]DocId)
	l.Unlock()
}
