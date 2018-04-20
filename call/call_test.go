package call

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
	"testing"
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
