package call

import (
	"gopkg.in/jarcoal/httpmock.v1"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestMakeCallsWhenURLExists(test *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://www.foo.com/bar",
		httpmock.NewStringResponder(200, `[]`))

	params := []string{"http://www.foo.com/bar", "10"}
	call, _ := BuildCall(params, 1, 100)
	results := call.MakeIt()

	assert.Equal(test, 1, len(results))
	assert.Equal(test, 10, results[200])
}

func TestMakeCallsWhenURLDoesntExist(test *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://www.foo.com/bar",
		httpmock.NewStringResponder(404, `[]`))

	params := []string{"http://www.foo.com/bar", "10"}
	call, _ := BuildCall(params, 1, 100)
	results := call.MakeIt()

	assert.Equal(test, 1, len(results))
	assert.Equal(test, 0, results[200])
	assert.Equal(test, 10, results[404])
}

func TestMakeCallsReturnTheSameStatusCode(test *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://www.foo.com/bar",
		httpmock.NewStringResponder(200, `[]`))

	params := []string{"http://www.foo.com/bar", "100"}
	call, _ := BuildCall(params, 1, 10)
	results := call.MakeIt()

	assert.Equal(test, 1, len(results))
	assert.Equal(test, 100, results[200])
}