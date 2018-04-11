package call

import (
	"fmt"
	"strconv"
	"net/http"
	"errors"
	"net/url"
	"log"
)

// A Call represents the very basic structure to
// start calling some URL out. It carries all data
// needed to call-it operate on.
type Call struct {
	url        string // The endpoint to be tested
	attempts   int    // number of attempts
	concurrent int    // number of concurrent calls
}

// Parses all given arguments and transform them into a Call
func BuildCall(args []string, maxAttempts int, maxConcurrentCalls int) (call Call, err error) {
	isValid, err := validate(args)
	if isValid == false {
		return
	}

	callUrl := args[0]
	attempts, err := ParseAttempts(args, maxAttempts)
	if err != nil {
		return
	}

	concurrentCalls, err := ParseConcurrentCalls(args, maxConcurrentCalls)
	if err != nil {
		return
	}

	call = Call{callUrl, attempts, concurrentCalls}
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
func ParseAttempts(args []string, defaultAttempts int) (attempts int, err error) {
	if len(args) == 1 {
		attempts = defaultAttempts
		return
	}

	attempts, attemptsErr := strconv.Atoi(args[1])
	if err != attemptsErr || attempts == 0 {
		fmt.Println("Number of attempts invalid. Using default: " + strconv.Itoa(defaultAttempts))
		attempts = defaultAttempts
	}
	return
}

// Tries to parse maxConcurrentCalls number. If it wasn't possible, returns
// default concurrent calls
func ParseConcurrentCalls(args []string, defaultConcurrentCalls int) (calls int, err error) {
	if len(args) <= 2 {
		calls = defaultConcurrentCalls
		return
	}

	calls, concurrentErr := strconv.Atoi(args[2])
	if err != concurrentErr || calls == 0 {
		fmt.Println("Number of concurrent calls. Using default: " + strconv.Itoa(defaultConcurrentCalls))
		calls = defaultConcurrentCalls
	}
	return
}

// Make a call and return its results
func MakeA(call Call) (results map[int]int) {
	results = make(map[int]int)

	fmt.Print("\n")
	for call.attempts > 0 {
		response, err := http.Get(call.url)
		if err != nil {
			log.Fatal("Something got wrong ", err)
		}

		results[response.StatusCode]++
		call.attempts--
		fmt.Print(". ")

		if call.attempts%30 == 0 {
			fmt.Print("\n")
		}
	}

	return
}

// Print results formatted by Status
func PrintResults(results map[int]int) {
	fmt.Println("\n\nResults:")
	for k, v := range results {
		fmt.Printf("Status " + strconv.Itoa(k) + " - " + strconv.Itoa(v) + " times\n")
	}
}
