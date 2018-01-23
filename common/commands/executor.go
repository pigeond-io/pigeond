// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package commands

import(
	"errors"
	"fmt"
)

type CallbackExecutor struct {
	action string
	args   [][]byte
}

func MakeExecutor(r Request) Executor {
	if r.Ok() {
		return &CallbackExecutor{action: r.Action(), args: r.Args()}
	} else {
		return &CallbackExecutor{action: "ERROR", args: make([][]byte, 0)} 
	}
}

func (r *CallbackExecutor) Execute(registry RegistryReader) error {
	callback, ok := registry.Read(r.action)
	var err error
	if ok {
		ok = callback(r.args...)
		if !ok {
			err = errors.New(fmt.Sprintf("%q failed", r.action))
		}
	} else {
		err = errors.New(fmt.Sprintf("Unknown Command %q", r.action))
	}
	return err
}
