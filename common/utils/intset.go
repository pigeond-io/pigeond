// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package utils

import (
	"sync"
	"time"
)

var (
	defaultDirtyCacheDuration = 5 * time.Minute //Default duration for which dirty cache is valid
)

/*
  Thread-safe IntSet with eventually consistent members cache
  IntSet implements DocIds interface
  Members cache is made dirty only on removal
*/
type IntSet struct {
	Count              int                //Member count
	members            []int64            //Eventually consistent Members Cache
	index              map[int64]struct{} //Index keyed with members
	dirty              int64              //Timestamp when the members cache got dirty
	dirtyCacheDuration time.Duration      //Duration for which dirty members cache is valid
	lock               sync.RWMutex       //ReadWrite synchronization mutex
}

// Constructor to create intsets. It uses first parameter as the dirtyCacheDuration as pointer to time.Duration, if nil dirtyCacheDuration is set to default dirty interval i.e. defaultDirtyCacheDuration. Followed by variable number of int64 members
func MakeIntSet(dirtyCacheDuration *time.Duration, args ...int64) *IntSet {
	if dirtyCacheDuration == nil {
		dirtyCacheDuration = &defaultDirtyCacheDuration
	}
	intset := &IntSet{
		index:              make(map[int64]struct{}),
		members:            make([]int64, 0),
		dirtyCacheDuration: *dirtyCacheDuration,
	}
	for _, a := range args {
		intset.add(a)
	}
	return intset
}

func (s *IntSet) Clear() {
	l := &s.lock
	l.Lock()
	s.index = make(map[int64]struct{})
	s.members = make([]int64, 0)
	s.dirty = 0
	l.Unlock()
}

func (s *IntSet) Add(a int64) {
	l := &s.lock
	l.RLock()
	_, ok := s.index[a]
	if !ok {
		l.RUnlock()
		l.Lock()
		s.add(a)
		l.Unlock()
	} else {
		l.RUnlock()
	}
}

func (s *IntSet) Remove(a int64) {
	l := &s.lock
	l.RLock()
	_, ok := s.index[a]
	if ok {
		l.RUnlock()
		l.Lock()
		delete(s.index, a)
		s.Count--
		s.setDirty()
		l.Unlock()
	} else {
		l.RUnlock()
	}
}

func (s *IntSet) Contains(a int64) bool {
	l := &s.lock
	l.RLock()
	_, ok := s.index[a]
	l.RUnlock()
	return ok
}

// Members is eventual consistence
// On Remove we mark the intset dirty
// We wait for MaxDirtyInterval before we cleanup
func (s *IntSet) Members() []int64 {
	l := &s.lock
	l.RLock()
	members := s.members
	if s.shouldCleanUp() {
		l.RUnlock()
		l.Lock()
		members = make([]int64, s.Count)
		index := 0
		for item := range s.index {
			members[index] = item
			index++
		}
		s.dirty = 0
		s.members = members
		l.Unlock()
	} else {
		l.RUnlock()
	}
	return members
}

func (s *IntSet) shouldCleanUp() bool {
	return s.isDirty() && (time.Since(time.Unix(s.dirty, 0)) > s.dirtyCacheDuration)
}

func (s *IntSet) isDirty() bool {
	return s.dirty != 0
}

func (s *IntSet) setDirty() {
	if !s.isDirty() {
		s.dirty = time.Now().Unix()
	}
}

// Critical Section that adds a member to set
func (s *IntSet) add(a int64) {
	s.index[a] = struct{}{}
	s.members = append(s.members, a)
	s.Count++
}
