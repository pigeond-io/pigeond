// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package main

import (
  "flag"
  "github.com/pigeond-io/pigeond/log"
  "github.com/pigeond-io/pigeond/stats"
  "github.com/pigeond-io/pigeond/utils"
  "github.com/pigeond-io/pigeond/edge"
)

var (
  addr      = flag.String("addr", "localhost:8765", "edge server address")
  debugMode = flag.Bool("debug", true, "enable debug mode")
  logFile   = flag.String("log", "", "log file")
)

func main() {
  flag.Parse()
  utils.InitProcess("edge", func(name string) {
    log.Init(name, *logFile, *debugMode)
  })
  utils.OnProcessExit(func() {
    //close file descriptors
  })
  println(utils.GetHeader())
  go stats.Logger()
  edge.InitWsServer(*addr)
}