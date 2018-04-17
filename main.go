package main

import (
	"log"
	"os"
	"github.com/urfave/cli"
	"fmt"
	"github.com/pedrolopesme/call-it/call"
)

const (
	DefaultAttempts        = call.Attempts(10)
	DefaultConcurrentCalls = call.Attempts(10)
)

func main() {
	app := cli.NewApp()
	app.Name = "Call It"
	app.Description = "A simple program to benchmark URL responses across multiple requests"
	app.Usage = "call-it [url] [number of attempts]"
	app.Version = "0.0.1-beta"

	app.Action = func(c *cli.Context) error {
		callAttempt, err := call.BuildCall(c.Args(), DefaultAttempts, DefaultConcurrentCalls)
		if err != nil {
			fmt.Println("It was impossible to parse arguments")
			os.Exit(1)
		}

		results := callAttempt.MakeIt()
		call.PrintResults(results)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
