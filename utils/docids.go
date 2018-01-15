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

// DocIndex is a common interface that encapsulates different implementations of inverted index of 64-bit Integers to Collection of 64-bit Integers
type DocIndex interface {
  Add(keyId int64, docId int64)           //Links keyId to docId
  Remove(keyId int64, docId int64)        //Removes keyId to docId
  Contains(keyId int64, docId int64) bool //Checks whether docId and keyId are linked or not
  Values(keyId int64) []int64             //Returns all the docIds that are linked to keyId
  Clear()                                 //Clears index
}
