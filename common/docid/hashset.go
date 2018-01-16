// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package docid

import (
	"sync"
	"time"
)

var (
	defaultDirtyCacheDuration = 5 * time.Minute //Default duration for which dirty cache is valid
)

/*
  Eventually consistent thread-safe implementation of Set interface
  Note: while adding new members the cache is always consistent.
  Caches becomes stale (dirty) when a member is removed from the Set
*/
type HashSet struct {
	Count              int                   //Member count
	members            []DocId               //Members Cache
	index              map[string]DocId //Index keyed with members
	dirty              int64                 //Timestamp when the members cache got dirty
	dirtyCacheDuration time.Duration         //Duration for which dirty members cache is valid
	lock               sync.RWMutex          //ReadWrite synchronization mutex
}

/*
	Constructor to create HashSets.
	first parameter is mandatory and takes dirtyCacheDuration time.Duration object as reference, if nil dirtyCacheDuration is set to default dirty interval i.e. defaultDirtyCacheDuration.
	next you can pass on all the members of type docId you would initialize set with.
*/
func MakeHashSet(dirtyCacheDuration *time.Duration, members ...DocId) *HashSet {
	if dirtyCacheDuration == nil {
		dirtyCacheDuration = &defaultDirtyCacheDuration
	}
	capacity := len(members) + 256
	set := &HashSet{
		index:              make(map[string]DocId),
		members:            make([]DocId, 0, capacity),
		dirtyCacheDuration: *dirtyCacheDuration,
	}
	for _, member := range members {
		set.add(member)
	}
	return set
}

func (s *HashSet) Clear() {
	l := &s.lock
	l.Lock()
	s.index = make(map[string]DocId)
	s.members = make([]DocId, 0, 256)
	s.dirty = 0
	l.Unlock()
}

//Adds member to HashSet
func (s *HashSet) Add(a DocId) {
	l := &s.lock
	id := a.Id()
	l.RLock()
	_, ok := s.index[id]
	if !ok {
		if cap(s.members) == len(s.members) && s.isDirty() {
			//Anyways the slice will grow even if we don't rebuild members.
			//Using this opportunity to grow the slice and rebuild the members from the index
			l.RUnlock()
			s.RebuildMembers()
		} else {
			l.RUnlock()
		}
		// l.RUnlock()
		l.Lock()
		s.add(a)
		l.Unlock()
	} else {
		l.RUnlock()
	}
}

//Removes member from HashSet
func (s *HashSet) Remove(a DocId) {
	id := a.Id()
	l := &s.lock
	l.RLock()
	_, ok := s.index[id]
	if ok {
		l.RUnlock()
		l.Lock()
		delete(s.index, id)
		s.Count--
		s.setDirty()
		l.Unlock()
	} else {
		l.RUnlock()
	}
}

//Checks whether a DocId belongs to the Set
func (s *HashSet) Contains(a DocId) bool {
	id := a.Id()
	l := &s.lock
	l.RLock()
	_, ok := s.index[id]
	l.RUnlock()
	return ok
}

//Publishes members to channel
func (s *HashSet) PublishMembers(blockSize int, channel chan []DocId) {
	l := &s.lock
	l.RLock()
	if s.shouldRebuildCache() {
		l.RUnlock()
		s.RebuildMembers()
		l.RLock()
	}
	members := s.members
	size := len(s.members)
	l.RUnlock()
	if blockSize < 1 {
		blockSize = 1024
	}
	// initiate a goroutine to send members to channel
	go func() {
		for start, end := 0, 0; start < size; start = end {
			if start + blockSize > size {
				end = size
			} else {
				end = start + blockSize
			}
			channel <- members[start:end]
		}
		channel <- []DocId{} // To Notify Termination
	}()
}

// Rebuilds Members Cache
func (s *HashSet) RebuildMembers() {
	l := &s.lock
	l.Lock()
	members := make([]DocId, s.Count, s.Count << 1)
	index := 0
	for _, docId := range s.index {
		members[index] = docId
		index++
	}
	s.dirty = 0
	s.members = members
	l.Unlock()
}

// Critical Section that adds a member to set
func (s *HashSet) add(a DocId) {
	s.index[a.Id()] = a
	s.members = append(s.members, a)
	s.Count++
}

func (s *HashSet) shouldRebuildCache() bool {
	return s.isDirty() && (time.Since(time.Unix(s.dirty, 0)) > s.dirtyCacheDuration)
}

func (s *HashSet) isDirty() bool {
	return s.dirty != 0
}

func (s *HashSet) setDirty() {
	if !s.isDirty() {
		s.dirty = time.Now().Unix()
	}
}
