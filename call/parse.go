package call

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

var (
	// ErrInvalidArgumentsNumber is an error when the number of arguments are invalid
	ErrInvalidArgumentsNumber = errors.New("invalid number of arguments")

	// ErrInvalidURL is an error when the format of the url is invalid
	ErrInvalidURL = errors.New("invalid url format")

	// ErrMethodNotAllowed is an error with bad config method
	ErrMethodNotAllowed = errors.New("Method not allowed")

	// ErrEmptyName is an error with bad Config method
	ErrEmptyName = errors.New("Request Config name cannot be nil")
)

const (
	// AttemptsPosition in Args
	AttemptsPosition = 1

	// ConcurrentAttemptsPosition in Args
	ConcurrentAttemptsPosition = 2
)

// BuildCall parses all given arguments and transform them into a ConcurrentCall
func BuildCall(args []string, maxAttempts, maxConcurrentAttempts int) (call ConcurrentCall, err error) {
	var (
		callURL                      *url.URL
		attempts, concurrentAttempts int
	)

	isValid, err := validate(args)
	if !isValid {
		return
	}

	callURL, err = url.Parse(args[0])
	if err != nil {
		return
	}

	attempts, err = ParseIntArgument(args, AttemptsPosition, maxAttempts)
	if err != nil {
		return
	}

	concurrentAttempts, err = ParseIntArgument(args, ConcurrentAttemptsPosition, maxConcurrentAttempts)
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
		return false, ErrInvalidArgumentsNumber
	}

	_, err = url.ParseRequestURI(args[0])
	if err != nil {
		return false, ErrInvalidURL
	}

	return true, nil
}

// ParseIntArgument tries to parse an int argument. If it wasn't possible, returns
// default value
func ParseIntArgument(args []string, position int, defaultValue int) (val int, err error) {
	if len(args) <= position {
		val = defaultValue
		return
	}

	val, err = strconv.Atoi(args[position])
	if err != nil {
		fmt.Println("Argument invalid. Using default: " + strconv.Itoa(int(defaultValue)))
		val = defaultValue
	}
	return
}

// BuildCallsFromConfig parses a Config file and transforms the instructions into a list of ConcurrentCalls
func BuildCallsFromConfig() (calls []ConcurrentCall, err error) {
	callConfig, err := config()
	if err != nil {
		return
	}
	for _, c := range callConfig {
		if err = c.checkDefaults(); err != nil {
			return
		}
		url, errP := url.ParseRequestURI(c.URL)
		if errP != nil {
			return nil, errP
		}
		newCall := ConcurrentCall{
			URL:                url,
			Attempts:           c.Attempts,
			ConcurrentAttempts: c.ConcurrentAttempts,
			config:             c,
		}
		calls = append(calls, newCall)
	}
	return
}
