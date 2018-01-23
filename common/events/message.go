// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package events

type SliceMessage struct {
	slice []byte
}

func (s *SliceMessage) Body() []byte {
	return s.slice
}

func MakeSliceMessage(slice []byte) (s *SliceMessage) {
	return &SliceMessage{slice: slice}
}