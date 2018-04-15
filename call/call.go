package call

import (
	"fmt"
	"net/http"
	"log"
)

// Represents a valid URL, including the protocol and all query strings
type URL string

// Number of call Attempts. Each call attempt results in a HTTP Status
// Code that will be used to build the benchmark at the end.
type Attempts int

// A Result is a map containing all StatusCodes and
// the total of occurrences.
type Result map[int]int

// A Call should know how to execute itself, generating
// a Result from its execution
type Call interface {
	MakeIt() Result
}

// A ConcurrentCall represents the very basic structure to
// start calling some URL out. It carries all data
// needed to call-it operate on.
type ConcurrentCall struct {
	URL                URL      // The endpoint to be tested
	Attempts           Attempts // number of Attempts
	ConcurrentAttempts Attempts // number of concurrent Attempts
}

// Make a call and return its results
func (call *ConcurrentCall) MakeIt() (results Result) {
	results = make(map[int]int)

	for call.Attempts > 0 {
		concurrentAttempts := calcTheNumberOfConcurrentAttempts(*call)
		statusCodeChannel := getUrl(call.URL, concurrentAttempts)
		for statusCode := range statusCodeChannel {
			results[statusCode]++
		}
		call.Attempts -= concurrentAttempts
	}

	return
}

// It calculates the amount of concurrent calls to be executed,
// based on the attempts left. It ensures that the next round
// of concurrent calls will respect the attempts of a given call
func calcTheNumberOfConcurrentAttempts(call ConcurrentCall) (numberOfConcurrentAttempts Attempts) {
	numberOfConcurrentAttempts = call.ConcurrentAttempts
	if numberOfConcurrentAttempts > call.Attempts {
		numberOfConcurrentAttempts = call.Attempts
	}
	return
}

// This func calls an URL concurrently
func getUrl(url URL, concurrentAttempts Attempts) chan int {
	statusCode := make(chan int)
	done := make(chan bool)

	for i := 0; i < int(concurrentAttempts); i ++ {
		go func() {
			response, err := http.Get(string(url))
			if err != nil {
				log.Fatal("Something got wrong ", err)
			}
			fmt.Print(" . ")
			statusCode <- response.StatusCode
			done <- true
		}()
	}

	go func() {
		for i := 0; i < int(concurrentAttempts); i ++ {
			<-done
		}
		close(statusCode)
	}()

	return statusCode
}
