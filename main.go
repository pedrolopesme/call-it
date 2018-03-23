package main

import (
	"log"
	"os"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Call It"
	app.Usage = "A simple program to benchmark URL responses across multiple requests"
	app.Version = "0.0.1-beta"

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
