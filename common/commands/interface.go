// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package commands

type ActionCallback func(...[]byte) bool

type Request interface {
	Ok() bool
	Action() string
	Args() [][]byte
}

type Executor interface {
	Execute(registry RegistryReader) error
}

type RegistryReader interface {
	Read(actionName string) (ActionCallback, bool)
}

type RegistryWriter interface {
	Write(actionName string, actionCallback ActionCallback) bool
	Close() bool
}

type RequestExecutor interface {
	Request
	Executor
}

type Registry interface {
	RegistryReader
	RegistryWriter
}
