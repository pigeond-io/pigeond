// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package docid

import (
	"runtime"
)

type SlicePublisher struct {
	slice []DocId
}

func MakeSlicePublisher(slice []DocId) (s *SlicePublisher) {
	sPublisher := &SlicePublisher{
		slice: slice,
	}
	return sPublisher
}

//Emit slice elements in block of blockSize to channel
func (s *SlicePublisher) Emit(channel chan []DocId, blockSize int) {
	if blockSize < 1 {
		blockSize = 1024 //default blockSize
	}
	// Initiate a goroutine
	go func() {
		for start, end, size := 0, 0, len(s.slice); start < size; start = end {
			if start+blockSize > size {
				end = size
			} else {
				end = start + blockSize
			}
			channel <- s.slice[start:end]
			// Allowing other goroutines to run
			runtime.Gosched()
		}
		// Sending Empty Slice for Termination
		channel <- []DocId{}
	}()
}
