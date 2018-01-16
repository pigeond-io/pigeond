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
)

func main() {
	println(utils.GetHeader())

	app := cli.NewApp()
	app.Name = "pigeond"
	app.Usage = "Start services"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "service",
			Value: "hub",
			Usage: "service name edge | hub | origin | data_store",
		},
	}

	app.Action = func(c *cli.Context) error {
		service := c.String("service")

		switch service {
		case "edge":
			edge.Init()
			break
		default:
			log.Error("Invalid service name")
			return errors.New("invalid service name")
		}
		return nil
	}

	app.Run(os.Args)
}
