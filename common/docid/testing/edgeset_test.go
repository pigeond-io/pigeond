// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package testing_test

import (
	"github.com/pigeond-io/pigeond/common/docid"
	. "github.com/pigeond-io/pigeond/common/docid/testing"
	"math/rand"
	"testing"
	// "time"
)

func BenchmarkHashEdgeSetAdd(b *testing.B) {
	edgeset := docid.MakeHashEdgeSet()
	nodes := 700
	slice := RandStrIds(nodes)
	var n int64
	var a, c docid.DocId
	for n = 0; n < int64(b.N); n++ {
		a = slice[rand.Intn(nodes)]
		c = slice[rand.Intn(nodes)]
		edgeset.Add(a, c)
	}
}

func TestHashEdgeSetSelfLoop(t *testing.T) {
	edgeset := docid.MakeHashEdgeSet()
	tt := TestEdgeSet(t, edgeset)
	a := ID(1)
	edgeset.Add(a, a)
	_ = tt.
		ShouldContain(a, a)
	edgeset.Remove(a, a)
	_ = tt.
		ShouldNotContain(a, a)
}

func TestHashEdgeSetDuplicates(t *testing.T) {
	edgeset := docid.MakeHashEdgeSet()
	tt := TestEdgeSet(t, edgeset)
	a := ID(1)
	b := ID(2)
	edgeset.Add(a, b)
	_ = tt.
		ShouldContain(a, b)
	edgeset.Add(a, b)
	_ = tt.
		ShouldContain(a, b)
	edgeset.Remove(a, b)
	_ = tt.
		ShouldNotContain(a, b)
}

func TestHashEdgeSet(t *testing.T) {
	edgeset := docid.MakeHashEdgeSet()
	tt := TestEdgeSet(t, edgeset)

	a := ID(1)
	b := ID(2)
	c := ID(3)
	d := ID(4)
	edgeset.Add(a, b)
	edgeset.Add(a, c)
	edgeset.Add(b, a)
	edgeset.Add(d, b)
	
	// Forward Edges Should be Contained
	_ = tt.
		ShouldContain(a, b).
		ShouldContain(a, c).
		ShouldContain(b, a).
		ShouldContain(d, b)
	// Backward Edges Should Not be Contained
	_ = tt.
		ShouldNotContain(c,a).
		ShouldNotContain(b,d).
		ShouldNotContain(c,d)

	// Sources Publisher
	publisher := edgeset.Sources(b)
	channel := make(chan []docid.DocId)
	index := 0
	publisher.Emit(channel, 0)
	for slice := range channel {
		if len(slice) == 0 {
			break
		}
		for _, docid := range slice {
			tt.ShouldContain(docid, b)
			index++
		}
	}
	_ = tt.
		CountShouldBe(2, index)
	close(channel)

	// Target Publisher
	publisher = edgeset.Targets(a)
	channel = make(chan []docid.DocId)
	publisher.Emit(channel, 0)
	index = 0
	for slice := range channel {
		if len(slice) == 0 {
			break
		}
		for _, docid := range slice {
			tt.ShouldContain(a, docid)
			index++
		}
	}
	_ = tt.CountShouldBe(2, index)
	close(channel)

	// Remove
	edgeset.Remove(d, b)
	tt.ShouldNotContain(d, b)

	edgeset.RemoveSource(a)
	tt.
		ShouldNotContain(a,b).
		ShouldNotContain(a,c).
		ShouldContain(b, a)
}
