// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package utils_test

import (
  "github.com/pigeond-io/pigeond/common/utils"
  "testing"
)

func BenchmarkIntSetAdd(b *testing.B) {
  // run the Fib function b.N times
  intset := utils.MakeIntSet()
  var n int64
  for n = 0; n < int64(b.N); n++ {
    intset.Add(n)
    intset.Remove(n >> 2)
  }
}

func BenchmarkIntSetAdd3Threads(b *testing.B) {
  // run the Fib function b.N times
  intset := utils.MakeIntSet()
  done := make(chan bool)
  worker := func(i int64) {
    var n int64
    for n = i; n < int64(b.N); n = n + 3 {
      intset.Add(n)
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

func TestIntSet(t *testing.T) {
  utils.MaxDirtyInterval = 0
  intset := utils.MakeIntSet()
  count := 0
  intset.Add(1)
  count++
  if intset.Count != count {
    t.Errorf("Count should be %d but is %d", count, intset.Count)
  }
  intset.Add(2)
  count++
  if intset.Count != count {
    t.Errorf("Count should be %d but is %d", count, intset.Count)
  }
  // TestDuplicate
  intset.Add(2)
  if intset.Count != 2 {
    t.Errorf("Count should be %d but is %d", count, intset.Count)
  }
  intset.Add(3)
  count++
  intset.Add(4)
  count++
  intset.Add(5)
  count++
  intset.Add(7)
  count++
  intset.Remove(2)
  count--
  if intset.Count != count {
    t.Errorf("Count should be %d but is %d", count, intset.Count)
  }
  // TestNonExistence
  intset.Remove(8)
  if intset.Count != count {
    t.Errorf("Count should be %d but is %d", count, intset.Count)
  }
  //Membership
  var (
    members    = []int64{1, 3, 4, 5, 7}
    nonmembers = []int64{2, 6, 8}
  )
  for _, member := range members {
    if !intset.Contains(member) {
      t.Errorf("Should contain %d", member)
    }
  }
  for _, member := range nonmembers {
    if intset.Contains(member) {
      t.Errorf("Should not contain %d", member)
    }
  }
  for _, member := range intset.Members() {
    if !intset.Contains(member) {
      t.Errorf("Should contain %d", member)
    }
  }
}
