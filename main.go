package main

import (
	"log"
	"os"
	"github.com/urfave/cli"
	"github.com/pedrolopesme/call-it/call"
	"fmt"
)

const (
	DefaultAttempts        = 10
	DefaultConcurrentCalls = 10
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

		results := call.MakeA(callAttempt)
		call.PrintResults(results)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
