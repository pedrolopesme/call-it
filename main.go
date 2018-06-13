package main

import (
	"log"
	"os"

	"github.com/pedrolopesme/call-it/call"
	"github.com/urfave/cli"
)

const (
	defaultAttempts        = 10
	defaultConcurrentCalls = 10
)

func main() {
	app := cli.NewApp()
	app.Name = "Call It"
	app.Description = "A simple program to benchmark URL responses across multiple requests"
	app.Usage = "call-it [url] [number of attempts]"
	app.Version = "0.0.1-beta"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "c",
			Usage: "Pass a config file with all data you want in one or more requests",
		},
	}
	app.Action = func(c *cli.Context) (err error) {
		if c.Bool("c") {
			return buildCallFromConfig()
		}
		return buildCall(c)
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func buildCall(c *cli.Context) (err error) {
	callAttempt, err := call.BuildCall(c.Args(), defaultAttempts, defaultConcurrentCalls)
	if err != nil {
		return
	}
	results := callAttempt.MakeIt()
	call.PrintResults(results)
	return
}

func buildCallFromConfig() (err error) {
	calls, err := call.BuildCallsFromConfig()
	if err != nil {
		return
	}
	for _, c := range calls {
		result := c.MakeIt()
		call.PrintResults(result)
	}
	return
}
