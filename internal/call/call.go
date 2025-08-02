package call

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
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
	config             Config   // configs from file
	Attempts           int      // number of Attempts
	ConcurrentAttempts int      // number of concurrent Attempts
}

// A Result contains the info to be outputted at the end
// of the operation
type Result struct {
	URL            *url.URL                    // Endpoint tested
	status         map[int]StatusCodeBenchmark // status codes
	totalExecution float64                     // total execution time
	avgExecution   float64                     // average execution time
	minExecution   float64                     // min execution time
	maxExecution   float64                     // min execution time
}

// HTTPResponse status code and execution time
type HTTPResponse struct {
	status    int     // status codes
	err       error   // possible error
	execution float64 // total execution time
}

// StatusCodeBenchmark with total of occurrences, execution time
type StatusCodeBenchmark struct {
	total     int     // total
	execution float64 // total execution time
}

// MakeIt executes a call and return its results
func (call *ConcurrentCall) MakeIt() (result Result) {
	result = Result{
		URL:            call.URL,
		status:         make(map[int]StatusCodeBenchmark),
		totalExecution: 0,
		avgExecution:   0,
		minExecution:   0,
		maxExecution:   0}

	beginning := time.Now()
	if call.config.Name != "" {
		fmt.Println("Case: ", call.config.Name)
	}
	s := spinner.New(spinner.CharSets[31], 300*time.Millisecond)
	s.Prefix = "😎 "
	s.Suffix = " " + call.URL.String()
	s.Start()
	totalAttempts := call.Attempts
	for call.Attempts > 0 {
		concurrentAttempts := calcConcurrentAttempts(*call)
		responses := callURL(call.URL, concurrentAttempts, call.config)
		for _, response := range responses {
			statusCodeBenchmark := result.status[response.status]
			statusCodeBenchmark.total++
			statusCodeBenchmark.execution += response.execution
			result.status[response.status] = statusCodeBenchmark
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

// GetStatus returns the status codes map for external access
func (r *Result) GetStatus() map[int]StatusCodeBenchmark {
	return r.status
}

// GetTotalExecution returns the total execution time
func (r *Result) GetTotalExecution() float64 {
	return r.totalExecution
}

// GetAvgExecution returns the average execution time
func (r *Result) GetAvgExecution() float64 {
	return r.avgExecution
}

// GetMinExecution returns the minimum execution time
func (r *Result) GetMinExecution() float64 {
	return r.minExecution
}

// GetMaxExecution returns the maximum execution time
func (r *Result) GetMaxExecution() float64 {
	return r.maxExecution
}

// GetTotal returns the total count for a status code benchmark
func (s *StatusCodeBenchmark) GetTotal() int {
	return s.total
}

// GetExecution returns the execution time for a status code benchmark
func (s *StatusCodeBenchmark) GetExecution() float64 {
	return s.execution
}

// It calculates the amount of concurrent calls to be executed,
// based on the attempts left. It ensures that the next round
// of concurrent calls will respect the attempts left of a given call
func calcConcurrentAttempts(call ConcurrentCall) (numberOfConcurrentAttempts int) {
	numberOfConcurrentAttempts = call.ConcurrentAttempts
	if numberOfConcurrentAttempts > call.Attempts {
		numberOfConcurrentAttempts = call.Attempts
	}
	return
}

// This func calls an URL concurrently
func callURL(callerURL *url.URL, concurrentAttempts int, config Config) (responses []HTTPResponse) {
	urlResponse := make(chan HTTPResponse)
	var wg sync.WaitGroup
	wg.Add(concurrentAttempts)
	for i := 0; i < int(concurrentAttempts); i++ {
		go func() {
			defer wg.Done()
			beginning := time.Now()
			req, err := buildRequest(callerURL.String(), config)
			if err != nil {
				log.Fatalf("Something got wrong: %v", err)
			}
			client := http.DefaultClient
			response, err := client.Do(req)
			executionSecs := time.Since(beginning).Seconds()
			if err != nil {
				urlResponse <- HTTPResponse{
					err:       err,
					execution: executionSecs,
					status:    http.StatusRequestTimeout,
				}
				return
			}
			urlResponse <- HTTPResponse{
				err:       err,
				execution: executionSecs,
				status:    response.StatusCode,
			}
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < int(concurrentAttempts); i++ {
			responses = append(responses, <-urlResponse)
		}
		close(urlResponse)
	}()
	wg.Wait()
	return
}

func buildRequest(baseURL string, config Config) (req *http.Request, err error) {
	if ok := reflect.DeepEqual(config, Config{}); ok {
		return http.NewRequest(http.MethodGet, baseURL, nil)
	}
	method := config.Method
	URL := config.URL
	body := strings.NewReader(config.Body)
	req, err = http.NewRequest(method, URL, body)
	if err != nil {
		return
	}
	req.Header = config.Header
	return
}
