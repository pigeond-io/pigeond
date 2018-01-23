// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package actions

import (
	"github.com/pigeond-io/pigeond/common/commands"
	"github.com/pigeond-io/pigeond/common/events"
)

func OnUnsubscribe(subscriber events.Subscriber) commands.ActionCallback {
	return func(args ...[]byte) bool {
		ok := true
		for arg := range args {
			ok = ok && subscriber.Unsubscribe(string(arg))
		}
		return ok
	}
}
