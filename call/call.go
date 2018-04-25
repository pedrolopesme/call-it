package call

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/briandowns/spinner"
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
	URL                *url.URL // The endpoint to be tested
	Attempts           int      // number of Attempts
	ConcurrentAttempts int      // number of concurrent Attempts
}

// A Result contains the info to be outputted at the end
// of the operation
type Result struct {
	status         map[int]int // status codes and the total of occurrences
	totalExecution float64     // total execution time
	avgExecution   float64     // average execution time
	minExecution   float64     // min execution time
	maxExecution   float64     // min execution time
}

// Get URL Response
type GetUrlResponse struct {
	status    int     // status codes and the total of occurrences
	execution float64 // total execution time
}

// Make a call and return its results
func (call *ConcurrentCall) MakeIt() (result Result) {
	result = Result{status: make(map[int]int),
		totalExecution: 0,
		avgExecution: 0,
		minExecution: 0,
		maxExecution: 0}

	beginning := time.Now()
	s := spinner.New(spinner.CharSets[31], 300*time.Millisecond)
	s.Prefix = "you "
	s.Suffix = " " + call.URL.String()
	s.Start()

	totalAttempts := call.Attempts
	for call.Attempts > 0 {
		concurrentAttempts := calcTheNumberOfConcurrentAttempts(*call)
		responses := getURL(call.URL, concurrentAttempts)
		for response := range responses {
			result.status[response.status]++

			if result.minExecution == 0 || result.minExecution > response.execution {
				result.minExecution = response.execution
			}

			if result.maxExecution == 0 || result.maxExecution < response.execution {
				result.maxExecution = response.execution
			}
		}
		call.Attempts -= concurrentAttempts
	}
	s.Stop()
	result.totalExecution = time.Since(beginning).Seconds()
	result.avgExecution = result.totalExecution / float64(totalAttempts)
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
func getURL(callerURL *url.URL, concurrentAttempts int) chan GetUrlResponse {
	urlResponse := make(chan GetUrlResponse)
	done := make(chan bool)

	for i := 0; i < int(concurrentAttempts); i++ {
		go func() {
			beginning := time.Now()
			response, err := http.Get(callerURL.String())
			if err != nil {
				log.Fatal("Something got wrong ", err)
			}
			executionSecs := time.Since(beginning).Seconds()
			urlResponse <- GetUrlResponse{ status: response.StatusCode, execution: executionSecs }
			done <- true
		}()
	}

	go func() {
		for i := 0; i < int(concurrentAttempts); i++ {
			<-done
		}
		close(urlResponse)
	}()

	return urlResponse
}
