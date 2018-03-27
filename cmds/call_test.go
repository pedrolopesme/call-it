package cmds

import (
	"testing"
	. "gopkg.in/check.v1"
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

func (s *MySuite) TestBuildCallWithoutAttempts(c *C) {
	params := []string{"http://www.dummy.com"}
	call, _:= BuildCall(params, 50)
	c.Assert(call.attempts, Equals, 50)
}

func (s *MySuite) TestBuildCallWithInValidAttemptsFormat(c *C) {
	params := []string{"http://www.dummy.com", "dummyAttempts"}
	call, _:= BuildCall(params, 50)
	c.Assert(call.attempts, Equals, 50)
}
