// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package utils

import (
  "fmt"
  "github.com/pigeond-io/pigeond/common/log"
  "os"
  "os/signal"
  "syscall"
)

func InitProcess(processName string, initClosure func(string)) {
  initClosure(fmt.Sprintf("%s-%d", processName, os.Getpid()))
}

func OnProcessExit(exitClosure func()) {
  sigs := make(chan os.Signal, 1)
  signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
  go func() {
    <-sigs
    log.Debug("Exiting...")
    exitClosure()
    os.Exit(0)
  }()
}
