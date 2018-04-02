package call

import (
	"testing"
	. "gopkg.in/check.v1"
	"gopkg.in/jarcoal/httpmock.v1"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestBuildCallWithValidParams(c *C) {
	params := []string{"http://www.dummy.com", "30"}
	call, _ := BuildCall(params, 50)

	c.Assert(call.url, Equals, params[0])
	c.Assert(call.attempts, Equals, 30)
}

func (s *MySuite) TestBuildCallWithoutParams(c *C) {
	var params []string
	_, err := BuildCall(params, 50)
	c.Assert(err, NotNil)
}

func (s *MySuite) TestBuildCallWithInValidUrl(c *C) {
	params := []string{"dummy", "30"}
	_, err := BuildCall(params, 50)
	c.Assert(err, NotNil)
}

func (s *MySuite) TestParseAttempts(c *C) {
	params := []string{"http://www.dummy.com", "10"}
	attempts, _:= ParseAttempts(params, 50)
	c.Assert(attempts, Equals, 10)
}

func (s *MySuite) TestParseAttemptsWithoutAttempts(c *C) {
	params := []string{"http://www.dummy.com"}
	attempts, _:= ParseAttempts(params, 50)
	c.Assert(attempts, Equals, 50)
}

func (s *MySuite) TestParseAttemptsInValidAttemptsFormat(c *C) {
	params := []string{"http://www.dummy.com", "dummyAttempts"}
	attempts, _:= ParseAttempts(params, 50)
	c.Assert(attempts, Equals, 50)
}

func (s *MySuite) TestMakeCallsWhenURLExists(c *C) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://www.foo.com/bar",
		httpmock.NewStringResponder(200, `[]`))

	params := []string{"http://www.foo.com/bar", "10"}
	call, _ := BuildCall(params, 1)
	results := MakeA(call)

	c.Assert(len(results), Equals, 1)
	c.Assert(results[200], Equals, 10)
}

func (s *MySuite) TestMakeCallsWhenURLDoesntExist(c *C) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://www.foo.com/bar",
		httpmock.NewStringResponder(404, `[]`))

	params := []string{"http://www.foo.com/bar", "10"}
	call, _ := BuildCall(params, 1)
	results := MakeA(call)

	c.Assert(len(results), Equals, 1)
	c.Assert(results[200], Equals, 0)
	c.Assert(results[404], Equals, 10)
}

func (s *MySuite) TestMakeCallsReturnTheSameStatusCode(c *C) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://www.foo.com/bar",
		httpmock.NewStringResponder(200, `[]`))

	params := []string{"http://www.foo.com/bar", "100"}
	call, _ := BuildCall(params, 1)
	results := MakeA(call)

	c.Assert(len(results), Equals, 1)
	c.Assert(results[200], Equals, 100)
}

func (s *MySuite) TestMakeCallsReturnMultipleStatusCodes(c *C) {
	c.Fail()
}