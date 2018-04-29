package call

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
)

func TestMakeCallsWhenURLExists(test *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://www.foo.com/bar",
		httpmock.NewStringResponder(200, `[]`))

	params := []string{"http://www.foo.com/bar", "10"}
	call, _ := BuildCall(params, 1, 100)
	result := call.MakeIt()

	assert.Equal(test, 1, len(result.status))
	assert.Equal(test, 10, result.status[200])
}

func TestMakeCallsWhenURLDoesntExist(test *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://www.foo.com/bar",
		httpmock.NewStringResponder(404, `[]`))

	params := []string{"http://www.foo.com/bar", "10"}
	call, _ := BuildCall(params, 1, 100)
	result := call.MakeIt()

	assert.Equal(test, 1, len(result.status))
	assert.Equal(test, 0, result.status[200])
	assert.Equal(test, 10, result.status[404])
}

func TestMakeCallsReturnTheSameStatusCode(test *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://www.foo.com/bar",
		httpmock.NewStringResponder(200, `[]`))

	params := []string{"http://www.foo.com/bar", "100"}
	call, _ := BuildCall(params, 1, 10)
	result := call.MakeIt()

	assert.Equal(test, 1, len(result.status))
	assert.Equal(test, 100, result.status[200])
}

func TestCalcConcurrentAttemptsWhenThereAreEnoughAttemptsLeft(test *testing.T) {
	urlAddress, _ := url.Parse("http://www.a.com")
	call := ConcurrentCall{URL: urlAddress, Attempts: 100, ConcurrentAttempts: 10}

	assert.Equal(test, 10, calcConcurrentAttempts(call))
}

func TestCalcConcurrentAttemptsWhenThereAreNotEnoughAttemptsLeft(test *testing.T) {
	urlAddress, _ := url.Parse("http://www.a.com")
	call := ConcurrentCall{URL: urlAddress, Attempts: 10, ConcurrentAttempts: 100}

	assert.Equal(test, 10, calcConcurrentAttempts(call))
}

func TestCalcConcurrentAttemptsWhenAttemptsLeftIsEqualToConcurrentAttempts(test *testing.T) {
	urlAddress, _ := url.Parse("http://www.a.com")
	call := ConcurrentCall{URL: urlAddress, Attempts: 10, ConcurrentAttempts: 10}

	assert.Equal(test, 10, calcConcurrentAttempts(call))
}

func TestGetUrlChannel(test *testing.T) {
	urlAddress := "http://www.foo.com/bar"

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", urlAddress,
		httpmock.NewStringResponder(200, `[]`))

	parsedURL, _ := url.Parse(urlAddress)
	statusCodeChannel := getURL(parsedURL, 50)

	reponsesCounter := 0
	for response := range statusCodeChannel {
		reponsesCounter++
		assert.Equal(test, 200, response.status)
	}
	assert.Equal(test, 50, reponsesCounter)
}
