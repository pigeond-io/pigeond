package main

import (
	. "github.com/pigeond-io/pigeond/core"
	"github.com/pigeond-io/pigeond/edge"
	"os"
	"gopkg.in/urfave/cli.v1"
	"errors"
)

func main()  {
	app := cli.NewApp()
	app.Name = "pigeond"
	app.Usage = "Start services"

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "service",
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
			Error.Println("Invalid service name")
			return errors.New("invalid service name")
		}
		return nil
	}

	app.Run(os.Args)
}