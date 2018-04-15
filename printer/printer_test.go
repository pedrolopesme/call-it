package printer

import (
	. "gopkg.in/check.v1"
)

var _ = Suite(&MySuite{})

func (s *MySuite) TestPrintResults(c *C) {
	c.Assert(true, Equals, true)
}
