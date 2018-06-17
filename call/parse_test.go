package call

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildCallWithValidParams(test *testing.T) {
	params := []string{"http://www.dummy.com", "30"}
	cl, err := BuildCall(params, 50, 10)
	assert.Nil(test, err)
	assert.Equal(test, params[0], cl.URL.String())
	assert.Equal(test, 30, cl.Attempts)
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
	attempts, err := ParseIntArgument(params, AttemptsPosition, 50)
	assert.Nil(test, err)
	assert.Equal(test, 10, attempts)
}

func TestParseAttemptsWithoutAttempts(test *testing.T) {
	params := []string{"http://www.dummy.com"}
	attempts, err := ParseIntArgument(params, AttemptsPosition, 50)
	assert.Nil(test, err)
	assert.Equal(test, 50, attempts)
}

func TestParseAttemptsWithInvalidFormat(test *testing.T) {
	params := []string{"http://www.dummy.com", "dummyAttempts"}
	attempts, err := ParseIntArgument(params, AttemptsPosition, 50)
	assert.NotNil(test, err)
	assert.Equal(test, 50, attempts)
}

func TestBuildCallWithConcurrentCalls(test *testing.T) {
	params := []string{"http://www.dummy.com", "30", "50"}
	call, err := BuildCall(params, 2, 100)
	assert.Nil(test, err)
	assert.Equal(test, params[0], call.URL.String())
	assert.Equal(test, 30, call.Attempts)
	assert.Equal(test, 50, call.ConcurrentAttempts)
}

func TestParseConcurrentCallsWithoutConcurrentParameter(test *testing.T) {
	params := []string{"http://www.dummy.com"}
	concurrentCalls, err := ParseIntArgument(params, ConcurrentAttemptsPosition, 100)
	assert.Nil(test, err)
	assert.Equal(test, 100, concurrentCalls)
}

func TestParseConcurrentCallsWithInValidFormat(test *testing.T) {
	params := []string{"http://www.dummy.com", "10", "dummyConcurrentCalls"}
	concurrentCalls, err := ParseIntArgument(params, ConcurrentAttemptsPosition, 100)
	assert.NotNil(test, err)
	assert.Equal(test, 100, concurrentCalls)
}

func TestArgumentsValidationWhenNoArgumentsAreSupplied(test *testing.T) {
	params := []string{}
	isValid, err := validate(params)
	assert.Equal(test, ErrInvalidArgumentsNumber, err)
	assert.False(test, isValid)
}

func TestURLValidationWhenUrlIsInvalid(test *testing.T) {
	params := []string{"invalidurl"}
	isValid, err := validate(params)
	assert.NotNil(test, err)
	assert.False(test, isValid)
}
