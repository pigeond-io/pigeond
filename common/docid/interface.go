// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package docid

var (
	EOF = &Nil{} //End of file. Empty DocId
)

//DocId defines an interface that has a doc id
type DocId interface {
	DocId() string
}

// Publisher defines an interface that can emits slice of DocId to a channel.
type Publisher interface {
	// Note: It should be terminated by an empty slice
	Emit(channel chan []DocId, inBatchSize int)
}

// ImmutableMap defines interface of read only Map
type ImmutableMap interface {
	//Retrieves the member from the map
	Get(id string) (DocId, error)
	//Checks whether entry by id exists or not
	Contains(id string) bool
}

// Map is a common interace that encapsulates a one-to-one map of Id to DocId
type Map interface {
	ImmutableMap
	//Adds a may entry from member.Id() to member in the map
	Add(member DocId) error
	//Removes an entry by id from the map
	Remove(id string) error
}

// ImmutableSet defines interface of read only sets of DocId
type ImmutableSet interface {
	//Checks whether item is a member
	Contains(item DocId) bool
	//Returns a publisher that emits members in the set
	Members() Publisher
}

// Set is a common interface that encapsulates different implementations for the collection of 64-bit Integers
type Set interface {
	ImmutableSet
	//Adds a member to collection
	Add(member DocId) error
	//Removes a member from collections
	Remove(member DocId) error
}

// Buffer is common interface that encapsulates different implementations for double buffering techniques
type DoubleBuffer interface {
	//Adds a member to buffer
	Add(member DocId) error
	//Return all the elements in buffer and clears the buffer
	Slice() []DocId
}

// ImmutableEdgeSet defines interface of read only EdgeSet
type ImmutableEdgeSet interface {
	//Checks whether source and target are linked or not
	Contains(source DocId, target DocId) bool
	//Returns a publisher that emits sources linked to target
	Sources(target DocId) Publisher
	//Returns a publisher that emits targets linked to source
	Targets(source DocId) Publisher
}

// EdgeSet is a common interface that encapsulates different implementations of many to many relationships. It represents as a collection of edges
type EdgeSet interface {
	ImmutableEdgeSet
	//Adds an edge from source to target
	Add(source DocId, target DocId) error
	//Reoves an edge from source to target
	Remove(source DocId, target DocId) error
	//Removes the source and all the links originating from source
	RemoveSource(source DocId) error
	//Removes the target and all the links terminating to target
	RemoveTarget(target DocId) error
}

type AddIndexEntryWriter interface {
	Add(key DocId, val DocId) error
}

type RemoveIndexEntryWriter interface {
	Remove(key DocId, val DocId) error
}

type AddIndexEntryWriterCallback func(AddIndexEntryWriter) error

type RemoveIndexEntryWriterCallback func(RemoveIndexEntryWriter) error

type ImmutableIndexMap interface {
	Query(indexTag int, key DocId) (Publisher, error)
	Add(indexTag int, callback AddIndexEntryWriterCallback) error
	Remove(indexTag int, callback RemoveIndexEntryWriterCallback) error
	RemoveKey(indexTag int, key DocId) error
	RemoveValue(indexTag int, val DocId) error
}
