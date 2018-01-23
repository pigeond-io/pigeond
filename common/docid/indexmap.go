// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package docid

import (
	"errors"
)

var (
	errorNotFound = errors.New("Index not found")
)

type HashImmutableIndexMap struct {
	indexmap map[int]EdgeSet
}

func MakeImmutableIndexMap(indexNames ...int) *HashImmutableIndexMap {
	indexmap := make(map[int]EdgeSet)
	for _, indexName := range indexNames {
		indexmap[indexName] = MakeHashEdgeSet()
	}
	return &HashImmutableIndexMap{indexmap: indexmap}
}

func (h *HashImmutableIndexMap) Query(indexName int, key DocId) (Publisher, error) {
	edgeset, ok := h.indexmap[indexName]
	if ok {
		return edgeset.Targets(key), nil
	}
	return nil, errorNotFound
}

func (h *HashImmutableIndexMap) Add(indexName int, callback AddIndexEntryWriterCallback) error {
	err := errorNotFound
	edgeset, ok := h.indexmap[indexName]
	if ok {
		err = callback(edgeset)
	}
	return err
}

func (h *HashImmutableIndexMap) Remove(indexName int, callback RemoveIndexEntryWriterCallback) error {
	err := errorNotFound
	edgeset, ok := h.indexmap[indexName]
	if ok {
		err = callback(edgeset)
	}
	return err
}

func (h *HashImmutableIndexMap) RemoveKey(indexName int, key DocId) error {
	err := errorNotFound
	edgeset, ok := h.indexmap[indexName]
	if ok {
		err = edgeset.RemoveSource(key)
	}
	return err
}

func (h *HashImmutableIndexMap) RemoveValue(indexName int, val DocId) error {
	err := errorNotFound
	edgeset, ok := h.indexmap[indexName]
	if ok {
		err = edgeset.RemoveTarget(val)
	}
	return err
}
