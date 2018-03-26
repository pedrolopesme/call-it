package main

import (
	"log"
	"os"
	"github.com/urfave/cli"
	"github.com/pedrolopesme/call-it/cmds"
)

const (
	DefaulAttempts = 10
)

func main() {
	app := cli.NewApp()
	app.Name = "Call It"
	app.Description= "A simple program to benchmark URL responses across multiple requests"
	app.Usage = "call-it [url] [number of attempts]"
	app.Version = "0.0.1-beta"

	app.Action = func(c *cli.Context) error {
		call := cmds.BuildCall(c.Args(), DefaulAttempts)
		cmds.MakeA(call)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
