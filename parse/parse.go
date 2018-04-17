package parse

import (
	"strconv"
	"fmt"
	"errors"
	"net/url"
	c "github.com/pedrolopesme/call-it/call"
)

// Parses all given arguments and transform them into a ConcurrentCall
func BuildCall(args []string, maxAttempts c.Attempts, maxConcurrentAttempts c.Attempts) (call c.ConcurrentCall, err error) {
	var callUrl c.URL
	var attempts c.Attempts
	var concurrentAttempts c.Attempts

	isValid, err := validate(args)
	if isValid == false {
		return
	}

	callUrl = c.URL(args[0])
	attempts, err = ParseAttempts(args, maxAttempts)
	if err != nil {
		return
	}

	concurrentAttempts, err = ParseConcurrentAttempts(args, maxConcurrentAttempts)
	if err != nil {
		return
	}

	call = c.ConcurrentCall{callUrl, attempts, concurrentAttempts}
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
func ParseAttempts(args []string, defaultAttempts c.Attempts) (attempts c.Attempts, err error) {
	if len(args) == 1 {
		attempts = defaultAttempts
		return
	}

	attemptsString, attemptsErr := strconv.Atoi(args[1])
	if err != attemptsErr || attempts == 0 {
		fmt.Println("Number of attempts invalid. Using default: " + strconv.Itoa(int(defaultAttempts)))
		attempts = defaultAttempts
	}
	attempts = c.Attempts(attemptsString)
	return
}

// Tries to parse the concurrent attempts number. If it wasn't possible, returns
// default concurrent attempts
func ParseConcurrentAttempts(args []string, defaultConcurrentAttempts c.Attempts) (attempts c.Attempts, err error) {
	if len(args) <= 2 {
		attempts = defaultConcurrentAttempts
		return
	}

	attemptsString, concurrentErr := strconv.Atoi(args[2])
	if err != concurrentErr || attempts == 0 {
		fmt.Println("Number of concurrent attempts. Using default: " + strconv.Itoa(int(defaultConcurrentAttempts)))
		attempts = defaultConcurrentAttempts
	}
	attempts = c.Attempts(attemptsString)
	return
}
