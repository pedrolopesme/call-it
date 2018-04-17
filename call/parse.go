package call

import (
	"strconv"
	"fmt"
	"errors"
	"net/url"
)

// Parses all given arguments and transform them into a ConcurrentCall
func BuildCall(args []string, maxAttempts Attempts, maxConcurrentAttempts Attempts) (call ConcurrentCall, err error) {
	var callUrl URL
	var attempts Attempts
	var concurrentAttempts Attempts

	isValid, err := validate(args)
	if isValid == false {
		return
	}

	callUrl = URL(args[0])
	attempts, err = ParseAttempts(args, maxAttempts)
	if err != nil {
		return
	}

	concurrentAttempts, err = ParseConcurrentAttempts(args, maxConcurrentAttempts)
	if err != nil {
		return
	}

	call = ConcurrentCall{callUrl, attempts, concurrentAttempts}
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
func ParseAttempts(args []string, defaultAttempts Attempts) (attempts Attempts, err error) {
	if len(args) == 1 {
		attempts = defaultAttempts
		return
	}

	attemptsString, attemptsErr := strconv.Atoi(args[1])
	if err != attemptsErr {
		fmt.Println("Number of attempts invalid. Using default: " + strconv.Itoa(int(defaultAttempts)))
		attempts = defaultAttempts
	} else {
		attempts = Attempts(attemptsString)
	}
	return
}

// Tries to parse the concurrent attempts number. If it wasn't possible, returns
// default concurrent attempts
func ParseConcurrentAttempts(args []string, defaultConcurrentAttempts Attempts) (attempts Attempts, err error) {
	if len(args) <= 2 {
		attempts = defaultConcurrentAttempts
		return
	}

	attemptsString, concurrentErr := strconv.Atoi(args[2])
	if err != concurrentErr {
		fmt.Println("Number of concurrent attempts. Using default: " + strconv.Itoa(int(defaultConcurrentAttempts)))
		attempts = defaultConcurrentAttempts
	} else {
		attempts = Attempts(attemptsString)
	}
	return
}
