package call

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// Parses all given arguments and transform them into a ConcurrentCall
func BuildCall(args []string, maxAttempts, maxConcurrentAttempts int) (call ConcurrentCall, err error) {
	var (
		callURL            *url.URL
		attempts           int
		concurrentAttempts int
	)

	isValid, err := validate(args)
	if isValid == false {
		return
	}

	callURL, err = url.Parse(args[0])
	if err != nil {
		return
	}
	attempts, err = ParseAttempts(args, maxAttempts)
	if err != nil {
		return
	}

	concurrentAttempts, err = ParseConcurrentAttempts(args, maxConcurrentAttempts)
	if err != nil {
		return
	}

	call = ConcurrentCall{
		URL:                callURL,
		Attempts:           attempts,
		ConcurrentAttempts: concurrentAttempts,
	}
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

	attempts, err = strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("Number of attempts invalid. Using default: " + strconv.Itoa(int(defaultAttempts)))
		attempts = defaultAttempts
	}
	return
}

// Tries to parse the concurrent attempts number. If it wasn't possible, returns
// default concurrent attempts
func ParseConcurrentAttempts(args []string, defaultConcurrentAttempts int) (attempts int, err error) {
	if len(args) <= 2 {
		attempts = defaultConcurrentAttempts
		return
	}

	attemptsString, concurrentErr := strconv.Atoi(args[2])
	if err != concurrentErr {
		fmt.Println("Number of concurrent attempts. Using default: " + strconv.Itoa(int(defaultConcurrentAttempts)))
		attempts = defaultConcurrentAttempts
	} else {
		attempts = attemptsString
	}
	return
}
