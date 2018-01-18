// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package docid

import (
	"runtime"
	"sync"
)

/*
  Thread-safe implementation of EdgeSet interface
*/
type HashEdgeSet struct {
	sourceEdges map[string]Set // Forward Edges
	targetEdges map[string]Set // Backward Edges
	lock        sync.RWMutex   //ReadWrite synchronization mutex
}

func MakeHashEdgeSet(members ...map[DocId]DocId) *HashEdgeSet {
	edges := &HashEdgeSet{
		sourceEdges: make(map[string]Set),
		targetEdges: make(map[string]Set),
	}
	for _, member := range members {
		for source, target := range member {
			sourceEdges, sourceExists := edges.sourceEdges[source.DocId()]
			targetEdges, targetExists := edges.targetEdges[target.DocId()]
			if !sourceExists {
				sourceEdges = edges.addSource(source)
			} 
			if !targetExists {
				targetEdges = edges.addTarget(target)
			}
			edges.addEdge(sourceEdges, target)
			edges.addEdge(targetEdges, source)
		}
	}
	return edges
}

//Adds an edge from sourceId to targetId
func (s *HashEdgeSet) Add(source DocId, target DocId) error {
	l := &s.lock
	l.RLock()
	sourceEdges, sourceExists := s.sourceEdges[source.DocId()]
	targetEdges, targetExists := s.targetEdges[target.DocId()]
	if !sourceExists || !targetExists {
		l.RUnlock()
		l.Lock()
		if !sourceExists {
			sourceEdges = s.addSource(source)
		} 
		if !targetExists {
			targetEdges = s.addTarget(target)
		}
		l.Unlock()
		l.RLock()
	}
	s.addEdge(sourceEdges, target)
	s.addEdge(targetEdges, source)
	l.RUnlock()
	return nil
}

//Removes an edge from source to target
func (s *HashEdgeSet) Remove(source DocId, target DocId) error {
	var err error
	l := &s.lock
	l.RLock()
	sourceEdges, sourceExists := s.sourceEdges[source.DocId()]
	targetEdges, targetExists := s.targetEdges[target.DocId()]
	if sourceExists && targetExists {
		err = sourceEdges.Remove(target)
		if err == nil {
			err = targetEdges.Remove(source)
		}
	}
	l.RUnlock()
	return err
}

//Checks whether edge between source and target
func (s *HashEdgeSet) Contains(source DocId, target DocId) bool {
	l := &s.lock
	l.RLock()
	sourceEdges, sourceExists := s.sourceEdges[source.DocId()]
	targetEdges, targetExists := s.targetEdges[target.DocId()]
	exists := sourceExists && targetExists
	if exists {
		exists = sourceEdges.Contains(target) && targetEdges.Contains(source)
	}
	l.RUnlock()
	return exists
}

//Publish to docIdSliceChannel in slices of size sliceBlockSize all the sources to target DocId. The stream is terminated with empty DocId slice
func (s *HashEdgeSet) Sources(target DocId) Publisher {
	l := &s.lock
	l.RLock()
	targetEdges, targetExists := s.targetEdges[target.DocId()]
	if targetExists {
		publisher := targetEdges.Members()
		l.RUnlock()
		return publisher
	} else {
		l.RUnlock()
		return MakeSlicePublisher([]DocId{})
	}
}

//Publish to docIdSliceChannel in slices of size sliceBlockSize all the targets from source DocId. The stream is terminated with empty DocId slice
func (s *HashEdgeSet) Targets(source DocId) Publisher {
	l := &s.lock
	l.RLock()
	sourceEdges, sourceExists := s.sourceEdges[source.DocId()]
	if sourceExists {
		publisher := sourceEdges.Members()
		l.RUnlock()
		return publisher
	} else {
		l.RUnlock()
		return MakeSlicePublisher([]DocId{})
	}
}

//Removes the source and all the links originating from source
func (s *HashEdgeSet) RemoveSource(source DocId) error {
	var err error
	l := &s.lock
	sourceId := source.DocId()
	l.RLock()
	sourceEdges, sourceExists := s.sourceEdges[source.DocId()]
	if sourceExists {
		channel := make(chan []DocId)
		publisher := sourceEdges.Members()
		l.RUnlock()
		l.Lock()
		delete(s.sourceEdges, sourceId)
		l.Unlock()
		publisher.Emit(channel, 0)
		for slice := range channel {
			if len(slice) == 0 {
				break
			}
			l.RLock()
			for _, target := range slice {
				targetEdges, targetExists := s.targetEdges[target.DocId()]
				if targetExists {
					targetEdges.Remove(source)
				}
			}
			l.RUnlock()
			runtime.Gosched()
		}
		close(channel)
	} else {
		l.RUnlock()
	}
	return err
}

//Removes the target and all the links terminating to target
func (s *HashEdgeSet) RemoveTarget(target DocId) error {
	var err error
	l := &s.lock
	targetId := target.DocId()
	l.RLock()
	targetEdges, targetExists := s.targetEdges[target.DocId()]
	if targetExists {
		channel := make(chan []DocId)
		publisher := targetEdges.Members()
		l.RUnlock()
		l.Lock()
		delete(s.targetEdges, targetId)
		l.Unlock()
		publisher.Emit(channel, 0)
		for slice := range channel {
			if len(slice) == 0 {
				break
			}
			l.RLock()
			for _, source := range slice {
				sourceEdges, sourceExists := s.sourceEdges[source.DocId()]
				if sourceExists {
					sourceEdges.Remove(target)
				}
			}
			l.RUnlock()
			runtime.Gosched()
		}
		close(channel)
	} else {
		l.RUnlock()
	}
	return err
}

func (s *HashEdgeSet) addSource(source DocId) Set {
	set := MakeHashSet(nil)
	s.sourceEdges[source.DocId()] = set
	return set 
}

func (s *HashEdgeSet) addTarget(target DocId) Set {
	set := MakeHashSet(nil)
	s.targetEdges[target.DocId()] = set
	return set
}

func (s *HashEdgeSet) addEdge(set Set, node DocId) {
	set.Add(node)
}
