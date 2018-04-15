package call

import (
	"testing"
	. "gopkg.in/check.v1"
	"github.com/pedrolopesme/call-it/parse"
	"gopkg.in/jarcoal/httpmock.v1"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestBuildCallWithValidParams(c *C) {
	params := []string{"http://www.dummy.com", "30"}
	cl, _ := parse.BuildCall(params, 50, 10)
	c.Assert(cl.URL, Equals, params[0])
	c.Assert(cl.Attempts, Equals, 30)
}

func (s *MySuite) TestBuildCallWithoutParams(c *C) {
	var params []string
	_, err := BuildCall(params, 50, 10)
	c.Assert(err, NotNil)
}

func (s *MySuite) TestBuildCallWithInValidUrl(c *C) {
	params := []string{"dummy", "30"}
	_, err := BuildCall(params, 50, 10)
	c.Assert(err, NotNil)
}

func (s *MySuite) TestParseAttempts(c *C) {
	params := []string{"http://www.dummy.com", "10"}
	attempts, _ := ParseAttempts(params, 50)
	c.Assert(attempts, Equals, 10)
}

func (s *MySuite) TestParseAttemptsWithoutAttempts(c *C) {
	params := []string{"http://www.dummy.com"}
	attempts, _ := ParseAttempts(params, 50)
	c.Assert(attempts, Equals, 50)
}

func (s *MySuite) TestParseAttemptsInValidAttemptsFormat(c *C) {
	params := []string{"http://www.dummy.com", "dummyAttempts"}
	attempts, _ := ParseAttempts(params, 50)
	c.Assert(attempts, Equals, 50)
}

func (s *MySuite) TestMakeCallsWhenURLExists(c *C) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://www.foo.com/bar",
		httpmock.NewStringResponder(200, `[]`))

	params := []string{"http://www.foo.com/bar", "10"}
	call, _ := BuildCall(params, 1, 10)
	results := call.MakeIt()

	c.Assert(len(results), Equals, 1)
	c.Assert(results[200], Equals, 10)
}

func (s *MySuite) TestMakeCallsWhenURLDoesntExist(c *C) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://www.foo.com/bar",
		httpmock.NewStringResponder(404, `[]`))

	params := []string{"http://www.foo.com/bar", "10"}
	call, _ := BuildCall(params, 1, 10)
	results := call.MakeIt()

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
	call, _ := BuildCall(params, 1, 10)
	results := call.MakeIt()

	c.Assert(len(results), Equals, 1)
	c.Assert(results[200], Equals, 100)
}

func (s *MySuite) TestBuildCallWithConcurrentCalls(c *C) {
	params := []string{"http://www.dummy.com", "30", "50"}
	call, _ := BuildCall(params, 50, 10)

	c.Assert(call.URL, Equals, params[0])
	c.Assert(call.Attempts, Equals, 30)
	c.Assert(call.ConcurrentAttempts, Equals, 50)
}

func (s *MySuite) TestParseConcurrentCallsWithoutConcurrentParameter(c *C) {
	params := []string{"http://www.dummy.com"}
	concurrentCalls, _ := ParseConcurrentAttempts(params, 100)
	c.Assert(concurrentCalls, Equals, 100)
}

func (s *MySuite) TestParseConcurrentCallsWithInValidFormat(c *C) {
	params := []string{"http://www.dummy.com", "10", "dummyConcurrentCalls"}
	concurrentCalls, _ := ParseConcurrentAttempts(params, 100)
	c.Assert(concurrentCalls, Equals, 100)
}
