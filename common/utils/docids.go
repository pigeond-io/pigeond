// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package utils

// DocIds is a common interface that encapsulates different implementations for the collection of 64-bit Integers
type DocIds interface {
  Add(member int64)         //Adds a member to collection
  Remove(member int64)      //Removes a member from collections
  Contains(item int64) bool //Checks whether item is a member
  Members() []int64         //Returns all the items
  Clear()                   //Clears the collection
}

// DocEdge is a common interface that encapsulates different implementations of many to many relationships. It represents as a collection of edges
type DocEdges interface {
  AddEdge(sourceId int64, targetId int64)           //Adds an edge from sourceId to targetId
  RemoveEdge(sourceId int64, targetId int64)        //Reoves an edge from sourceId to targetId
  RemoveSource(sourceId int64)                      //Removes the sourceId and all the links originating from sourceId
  RemoveTarget(targetId int64)                      //Removes the targetId and all the links terminating to targetId
  ContainsEdge(sourceId int64, targetId int64) bool //Checks whether edge between and keyId are linked or not
  SourceIds(targetId int64) []int64                 //Returns all the sourceIds where an edge exists to targetId
  TargetIds(sourceId int64) []int64                 //Returns all the targetIds where an edge exists from sourceId
  Clear()                                           //Clears
}
