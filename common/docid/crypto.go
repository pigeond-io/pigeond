// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package docid

import (
	"crypto/md5"
	"encoding/hex"
)

var (
	signature = []byte("7k86F#5$3KiG")
)

func MD5(bytes ...[]byte) string {
	hash := md5.New()
	for _, buf := range bytes {
		hash.Write(buf)
	}
	return hex.EncodeToString(hash.Sum(signature))
}