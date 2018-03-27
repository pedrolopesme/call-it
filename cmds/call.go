package cmds

import (
	"fmt"
	"strconv"
	"net/http"
	"errors"
	"net/url"
)

// A Call represents the very basic structure to
// start calling some URL out. It carries all data
// needed to call-it operate on.
type Call struct {
	url      string // The endpoint to be tested
	attempts int    // the number of attempts
}

// Parses all given arguments and transform them into a Call
func BuildCall(args []string, maxAttempts int) (call Call, err error) {
	isValid, err := validate(args)
	if isValid == false {
		return
	}

	callUrl := args[0]
	attempts, err := getAttempts(args, maxAttempts)
	if err != nil {
		return
	}

	call = Call{callUrl, attempts}
	return
}

// Checks if the given parameters are valid
func validate(args []string) (result bool, err error) {
	if args == nil || len(args) < 1 {
		return false, errors.New("invalidArguments")
	}

	_, err = url.ParseRequestURI(args[0])
	if err != nil {
		return false, errors.New("invalidUrl ")
	}

	return true, nil
}

// Tries to parse maxAttempts number. If it wasn't possible, returns
// default attempts
func getAttempts(args []string, defaultAttempts int) (attempts int, err error) {
	if(len(args) == 1) {
		attempts = 	defaultAttempts
		return
	}

	attempts, attempsErr := strconv.Atoi(args[1])
	if err != attempsErr || attempts == 0 {
		fmt.Println("Number of attemps invalid. Using default: " + strconv.Itoa(defaultAttempts))
		attempts = defaultAttempts
	}
	return
}

// Make a call
func MakeA(call Call) {
	results := make(map[int]int)

	fmt.Print("\n")
	for call.attempts > 0 {
		response, _ := http.Get(call.url)
		results[response.StatusCode]++
		call.attempts--
		fmt.Print(". ")

		if call.attempts%30 == 0 {
			fmt.Print("\n")
		}
	}

	fmt.Println("\n\nResults:")
	for k, v := range results {
		fmt.Printf("Status " + strconv.Itoa(k) + " - " + strconv.Itoa(v) + " times\n")
	}
}
