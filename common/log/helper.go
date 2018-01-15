// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package log

import (
  joonix "github.com/joonix/log"
  "github.com/sirupsen/logrus"
  "github.com/NebulousLabs/fastrand"
  "os"
  "path/filepath"
  "runtime"
  "time"
  "unsafe"
)

var logger = logrus.New()
var scriptName = "pigeond"
var debugMode = false
var workingDir = ""
var fieldNames = [...]string{"service", "method", "args", "@"}

func Init(ScriptName string, FilePath string, DebugMode bool) {
  workingDir, _ = os.Getwd()
  logger.Formatter = &joonix.FluentdFormatter{}
  scriptName = ScriptName
  if FilePath == "" {
    logger.Out = os.Stdout
  } else {
    file, err := os.OpenFile(FilePath, os.O_CREATE|os.O_WRONLY, 0666)
    if err == nil {
      logger.Out = file
    } else {
      logger.Out = os.Stdout
      Error("Failed to open file. Using default stdout.")
    }
  }
  if DebugMode {
    debugMode = true
    logger.SetLevel(logrus.DebugLevel)
  } else {
    debugMode = false
    logger.SetLevel(logrus.InfoLevel)
  }
}

// WithFields(service)
// WithFields(service, method)
// WithFields(service, method, args)
func WithFields(args ...interface{}) *logrus.Entry {
  fields := logrus.Fields{
    "script": scriptName,
  }
  if len(args) == 0 {
    if debugMode {
      _, file, line, _ := runtime.Caller(2)
      fields["!"], _ = filepath.Rel(workingDir, file)
      fields["#"] = line
    }
  } else {
    for i := range args {
      if i > 3 {
        break
      }
      fields[fieldNames[i]] = args[i]
    }
  }
  return logger.WithFields(fields)
}

func Instrument(start time.Time, onElapse func(string)) {
    if debugMode {
      // Every request instrument
      elapsed := time.Since(start)
      onElapse(elapsed.String())
    } else {
      // Randomly instrument
      b := fastrand.Bytes(1)
      r := *(*uint8)(unsafe.Pointer(&b[0]))
      if (r & ((r >> 4) << 4)) == r  {
        elapsed := time.Since(start)
        onElapse(elapsed.String())
      }
    }
}

func Debug(args ...interface{}) {
  WithFields().Debug(args)
}

func Info(args ...interface{}) {
  WithFields().Info(args)
}

func Error(args ...interface{}) {
  WithFields().Error(args)
}

func Fatal(args ...interface{}) {
  WithFields().Fatal(args)
}
