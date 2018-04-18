package call

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// A Call should know how to execute itself, generating
// a Result from its execution
type Call interface {
	MakeIt() map[int]int
}

// A ConcurrentCall represents the very basic structure to
// start calling some URL out. It carries all data
// needed to call-it operate on.
type ConcurrentCall struct {
	URL                url.URL // The endpoint to be tested
	Attempts           int     // number of Attempts
	ConcurrentAttempts int     // number of concurrent Attempts
}

// Make a call and return its results
func (call *ConcurrentCall) MakeIt() (results map[int]int) {
	results = make(map[int]int)

	for call.Attempts > 0 {
		concurrentAttempts := calcTheNumberOfConcurrentAttempts(*call)
		statusCodeChannel := getURL(call.URL, concurrentAttempts)
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
func calcTheNumberOfConcurrentAttempts(call ConcurrentCall) (numberOfConcurrentAttempts int) {
	numberOfConcurrentAttempts = call.ConcurrentAttempts
	if numberOfConcurrentAttempts > call.Attempts {
		numberOfConcurrentAttempts = call.Attempts
	}
	return
}

// This func calls an URL concurrently
func getURL(callerURL url.URL, concurrentAttempts int) chan int {
	statusCode := make(chan int)
	done := make(chan bool)

	for i := 0; i < int(concurrentAttempts); i++ {
		go func() {
			response, err := http.Get(callerURL.String())
			if err != nil {
				log.Fatal("Something got wrong ", err)
			}
			fmt.Print(" . ")
			statusCode <- response.StatusCode
			done <- true
		}()
	}

	go func() {
		for i := 0; i < int(concurrentAttempts); i++ {
			<-done
		}
		close(statusCode)
	}()

	return statusCode
}
