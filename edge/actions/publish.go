// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package actions

import (
	"github.com/pigeond-io/pigeond/common/commands"
	"github.com/pigeond-io/pigeond/common/events"
)

func OnPublish(publisher events.Publisher) commands.ActionCallback {
	return func(args ...[]byte) bool {
		topic := string(args[0])
		msgs := args[1:]
		evmsgs := make([]events.Message, 0, len(msgs))
		for _, msg := range msgs {
			evmsgs = append(evmsgs, events.MakeSliceMessage(msg))
		}
		return publisher.Publish(topic, evmsgs...)
	}
}
