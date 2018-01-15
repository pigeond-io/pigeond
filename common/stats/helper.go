// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package stats

import(
  "sync/atomic"
  "time"
  "fmt"
  "github.com/pigeond-io/pigeond/common/log"
  "runtime"
)

const (
  StatsTickInterval = 1*time.Second
)

var (
  served    int64
  live int64
  failed    int64
)

func IncrLive(){
  atomic.AddInt64(&live, 1)
}

func DecrLive(){
  atomic.AddInt64(&live, -1)
}

func IncrFailed(){
  atomic.AddInt64(&failed, 1)
}

func IncrServed(){
  atomic.AddInt64(&served, 1)
}

func Logger() {
  lastUpdate := ""
  for {
    time.Sleep(StatsTickInterval)
    currUpdate := fmt.Sprintf("goroutines = %d, served = %d, live = %d, failed = %d", runtime.NumGoroutine(), atomic.LoadInt64(&served), atomic.LoadInt64(&live), atomic.LoadInt64(&failed))
    if currUpdate != lastUpdate {
      log.WithFields("stats").Info(currUpdate)
      lastUpdate = currUpdate
    }
  }
}