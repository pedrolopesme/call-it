package call

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildCallWithValidParams(test *testing.T) {
	params := []string{"http://www.dummy.com", "30"}
	cl, _ := BuildCall(params, 50, 10)

	assert.Equal(test, URL(params[0]), cl.URL)
	assert.Equal(test, Attempts(30), cl.Attempts)
}

func TestBuildCallWithoutParams(test *testing.T) {
	var params []string
	_, err := BuildCall(params, 50, 10)

	assert.NotNil(test, err)
}

func TestBuildCallWithInValidUrl(test *testing.T) {
	params := []string{"dummy", "30"}
	_, err := BuildCall(params, 50, 10)
	assert.NotNil(test, err)
}

func TestParseAttempts(test *testing.T) {
	params := []string{"http://www.dummy.com", "10"}
	attempts, _ := ParseAttempts(params, 50)

	assert.Equal(test, Attempts(10), attempts)
}

func TestParseAttemptsWithoutAttempts(test *testing.T) {
	params := []string{"http://www.dummy.com"}
	attempts, _ := ParseAttempts(params, 50)
	assert.Equal(test, Attempts(50), attempts)
}

func TestParseAttemptsWithInvalidFormat(test *testing.T) {
	params := []string{"http://www.dummy.com", "dummyAttempts"}
	attempts, _ := ParseAttempts(params, 50)

	assert.Equal(test, Attempts(50), attempts)
}

func TestBuildCallWithConcurrentCalls(test *testing.T) {
	params := []string{"http://www.dummy.com", "30", "50"}
	call, _ := BuildCall(params, 100, 200)

	assert.Equal(test, URL(params[0]), call.URL)
	assert.Equal(test, Attempts(30), call.Attempts)
	assert.Equal(test, Attempts(50), call.ConcurrentAttempts)
}

func TestParseConcurrentCallsWithoutConcurrentParameter(test *testing.T) {
	params := []string{"http://www.dummy.com"}
	concurrentCalls, _ := ParseConcurrentAttempts(params, 100)
	assert.Equal(test, Attempts(100), concurrentCalls)
}

func TestParseConcurrentCallsWithInValidFormat(test *testing.T) {
	params := []string{"http://www.dummy.com", "10", "dummyConcurrentCalls"}
	concurrentCalls, _ := ParseConcurrentAttempts(params, 100)
	assert.Equal(test, Attempts(100), concurrentCalls)
}
