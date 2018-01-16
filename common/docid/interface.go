// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package docid

var (
	EOF = &Nil{} //End of file
)

type DocId interface {
	Id() string
}

// Map is a common interace that encapsulates a one-to-one map of Id to DocId
type Map interface {
	//Retrieves the member from the map
	Get(id string) (DocId, error)
	//Adds a may entry from member.Id() to member in the map
	Add(member DocId)
	//Removes an entry by id from the map
	Remove(id string)
	//Checks whether entry by id exists or not
	Contains(id string) bool
}

// Set is a common interface that encapsulates different implementations for the collection of 64-bit Integers
type Set interface {
	//Adds a member to collection
	Add(member DocId)
	//Removes a member from collections
	Remove(member DocId)
	//Checks whether item is a member
	Contains(item DocId) bool
	//Publish to docIdSliceChannel in slices of size sliceBlockSize all the members of the set. The stream is terminated with empty DocId slice
	PublishMembers(blockSize int, docIdSliceChannel chan []DocId)
}

// EdgeSet is a common interface that encapsulates different implementations of many to many relationships. It represents as a collection of edges
type EdgeSet interface {
	//Adds an edge from sourceId to targetId
	Add(source DocId, target DocId)
	//Reoves an edge from sourceId to targetId
	Remove(source DocId, target DocId)
	//Checks whether edge between and keyId are linked or not
	Contains(source DocId, target DocId) bool
	//Removes the sourceId and all the links originating from sourceId
	RemoveSource(source DocId)
	//Removes the targetId and all the links terminating to targetId
	RemoveTarget(target DocId)
	//Publish to docIdSliceChannel in slices of size sliceBlockSize all the sources to target DocId. The stream is terminated with empty DocId slice
	PublishSources(target DocId, sliceBlockSize int, docIdSliceChannel chan []DocId)
	//Publish to docIdSliceChannel in slices of size sliceBlockSize all the targets from source DocId. The stream is terminated with empty DocId slice
	PublishTargets(source DocId, sliceBlockSize int, docIdSliceChannel chan []DocId)
}
