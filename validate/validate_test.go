package validate

import (
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestIsJson(c *C) {
	c.Assert(IsJson(`{"name":"kain"}`), Equals, true)
	c.Assert(IsJson(`{}`), Equals, true)
	c.Assert(IsJson(`[]`), Equals, true)
	c.Assert(IsJson(`{"name":kkk}`), Equals, false)
	c.Assert(IsJson(`hello`), Equals, false)
}

func (s *MySuite) TestIsUuid(c *C) {
	c.Assert(IsUuid("asd"), Equals, false)
	c.Assert(IsUuid("1887ca4e-6675-48cc-a482-8dc4137155b4"), Equals, true)
}
