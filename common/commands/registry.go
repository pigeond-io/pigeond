// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package commands

/* This registry is not thread-safe.*/
type MapRegistry struct {
	actions map[string]ActionCallback
}

func MakeRegistry() Registry {
	return &MapRegistry{actions: make(map[string]ActionCallback)}
}

func (r *MapRegistry) Write(actionName string, onAction ActionCallback) bool {
	r.actions[actionName] = onAction
	return true
}

func (r *MapRegistry) Read(actionName string) (ActionCallback, bool) {
	val, ok := r.actions[actionName]
	return val, ok
}

func (r *MapRegistry) Close() bool {
	r.actions = nil
	return true
}
