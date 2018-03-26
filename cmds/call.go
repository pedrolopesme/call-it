package cmds

import (
	"fmt"
	"strconv"
	"net/http"
	"github.com/urfave/cli"
)

// A Call represents the very basic structure to
// start calling some URL out. It carries all data
// needed to call-it operate on.
type Call struct {
	url      string // The endpoint to be tested
	attempts int    // the number of attempts
}

// Parses all given arguments and transform them into a Call
func BuildCall(args cli.Args, maxAttempts int) (call Call){
	url := args.Get(0)
	attempts, err := strconv.Atoi(args.Get(1))
	if err != nil {
		fmt.Println("Number of attemps invalid. Using default: " + strconv.Itoa(maxAttempts))
		attempts = maxAttempts
	}
	call = Call{url, attempts}
	return
}

// Make a call
func MakeA(call Call){
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