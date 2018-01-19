// Copyright 2018 The PigeonD Authors. All rights reserved.
// Use of this source code is governed by a GNU AGPL v3.0
// license that can be found in the AGPL V3 LICENSE file.

package main

import (
	"errors"
	"github.com/pigeond-io/pigeond/common/log"
	"github.com/pigeond-io/pigeond/common/utils"
	"github.com/pigeond-io/pigeond/edge"
	"gopkg.in/urfave/cli.v1"
	"os"
	"github.com/pigeond-io/pigeond/common/stats"
	"bufio"
)

var flags = []cli.Flag{
	cli.StringFlag{
		Name:  "service",
		Value: "edge",
		Usage: "service name should be edge | hub | origin | data_store",
	},
	cli.StringFlag{
		Name:  "ws-address",
		Value: "localhost:8765",
		Usage: "websocket port",
	},
	cli.IntFlag{
		Name:  "ws-port",
		Value: 8001,
		Usage: "websocket port",
	},
	cli.IntFlag{
		Name:  "udp-port",
		Value: 8002,
		Usage: "udp port",
	},
	cli.IntFlag{
		Name:  "ws-buffer-size",
		Value: 2048,
		Usage: "websocket read buffer size",
	},
	cli.BoolFlag{
		Name:  "debug",
		Usage: "enable debug mode",
	},
	cli.StringFlag{
		Name:  "log",
		Value: "",
		Usage: "log file path",
	},
}

func main() {

	f := bufio.NewWriter(os.Stdout)
	f.Write([]byte(utils.GetHeader()))
	f.Flush()

	go stats.Logger()

	app := cli.NewApp()
	app.Name = "pigeond"
	app.Usage = "Start services"

	app.Flags = flags

	app.Action = func(c *cli.Context) error {

		//Set logging
		logFile := c.String("log")
		debugMode := c.Bool("debug")
		utils.InitProcess("edge", func(name string) {
			log.Init(name, logFile, debugMode)
		})
		utils.OnProcessExit(func() {
			//close file descriptors
		})

		service := c.String("service")

		switch service {
		case "edge":
			addr := c.String("ws-address")
			// wsPort := c.Int("wd-port")
			// udpPort := c.Int("udp-port")
			// wsBufferSize := c.Int("ws-buffer-size")
			// edge.Init(wsPort, udpPort, wsBufferSize)
			edge.InitWsServer(addr)
			break
		default:
			log.Error("Invalid service name")
			return errors.New("invalid service name")
		}
		return nil
	}

	app.Run(os.Args)
}
