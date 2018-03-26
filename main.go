package main

import (
	"log"
	"os"
	"github.com/urfave/cli"
	"fmt"
	"strconv"
	"net/http"
)

const (
	DefaulAttempts = 10
)

// A Call represents the very basic structure to
// start calling some URL out. It carries all data
// needed to call-it operate on.
type Call struct {
	url      string // The endpoint to be tested
	attempts int    // the number of attempts
}

// Parses all given arguments and transform them into a Call
func buildCall(args cli.Args) (call Call){
	url := args.Get(0)
	attempts, err := strconv.Atoi(args.Get(1))
	if err != nil {
		fmt.Println("Number of attemps invalid. Using default: " + strconv.Itoa(DefaulAttempts))
		attempts = DefaulAttempts
	}
	call = Call{url, attempts}
	return
}

// Make a call
func makeA(call Call){
	results := make(map[int]int)

	fmt.Print("\n")
	for call.attempts > 0{
		response, _ := http.Get(call.url)
		results[response.StatusCode]++
		call.attempts--
		fmt.Print(". ")

		if call.attempts % 30 == 0{
			fmt.Print("\n")
		}
	}

	fmt.Println("\n\nResults:")
	for k, v := range results {
		fmt.Printf("Status "+ strconv.Itoa(k) + " - " + strconv.Itoa(v) + " times\n")
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "Call It"
	app.Description= "A simple program to benchmark URL responses across multiple requests"
	app.Usage = "call-it [url] [number of attempts]"
	app.Version = "0.0.1-beta"

	app.Action = func(c *cli.Context) error {
		call := buildCall(c.Args())
		makeA(call)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
