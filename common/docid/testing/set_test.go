// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package testing_test

import (
	"github.com/pigeond-io/pigeond/common/docid"
	. "github.com/pigeond-io/pigeond/common/docid/testing"
	"testing"
	"time"
)

func BenchmarkMakeHashSet(b *testing.B) {
	var n int64
	for n = 0; n < int64(b.N); n++ {
		_ = docid.MakeHashSet(nil)
	}
}

func BenchmarkHashSetAdd(b *testing.B) {
	set := docid.MakeHashSet(nil)
	var n int64
	for n = 0; n < int64(b.N); n++ {
		set.Add(ID(n))
		set.Remove(ID(n >> 2))
	}
}

func BenchmarkHashSetAdd3Threads(b *testing.B) {
	set := docid.MakeHashSet(nil)
	done := make(chan bool)
	worker := func(i int64) {
		var n int64
		for n = i; n < int64(b.N); n = n + 3 {
			set.Add(ID(n))
		}
		done <- true
	}
	go worker(0)
	go worker(1)
	go worker(2)
	<-done
	<-done
	<-done
}

func TestHashSet(t *testing.T) {
	cacheDirtyValidity := time.Duration(0)
	set := docid.MakeHashSet(&cacheDirtyValidity)
	tt := TestSet(t, set)
	count := 0
	set.Add(ID(1))
	count++
	_ = tt.CountShouldBe(count, set.Count)
	set.Add(ID(2))
	count++
	_ = tt.CountShouldBe(count, set.Count)
	// TestDuplicate
	set.Add(ID(2))
	_ = tt.CountShouldBe(count, set.Count)
	set.Add(ID(3))
	count++
	set.Add(ID(4))
	count++
	set.Add(ID(5))
	count++
	set.Add(ID(7))
	count++
	set.Remove(ID(2))
	count--
	_ = tt.CountShouldBe(count, set.Count)
	// TestNonExistence
	set.Remove(ID(8))
	_ = tt.CountShouldBe(count, set.Count)
	//Membership
	set.RebuildMembers()
	publisher := set.Members()
	channel := make(chan []docid.DocId)
	publisher.Emit(channel, 0)
	for slice := range channel {
		if len(slice) == 0 {
			break
		}
		index := 0
		for _, docid := range slice {
			_ = tt.ShouldContain(docid)
			index++
		}
		_ = tt.CountShouldBe(set.Count, index)
	}
	var (
		nonmembers = []int64{2, 6, 8}
	)
	for _, member := range nonmembers {
		_ = tt.ShouldNotContain(ID(member))
	}
}
