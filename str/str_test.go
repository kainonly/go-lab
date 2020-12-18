package str

import (
	"github.com/go-playground/validator/v10"
	. "gopkg.in/check.v1"
	"testing"
)

var err error

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestRandom(c *C) {
	result := Random(8)
	if err = validator.New().Var(result, "alpha,len=8"); err != nil {
		c.Error(err)
	}
	c.Log(result)
}

func (s *MySuite) TestIsJson(c *C) {
	c.Assert(IsJson(`{"name":"kain"}`), Equals, true)
	c.Assert(IsJson(`{}`), Equals, true)
	c.Assert(IsJson(`[]`), Equals, true)
	c.Assert(IsJson(`{"name":kkk}`), Equals, false)
	c.Assert(IsJson(`hello`), Equals, false)
}

func (s *MySuite) TestIsUuid(c *C) {
	c.Assert(IsUuid("asd"), Equals, false)
	c.Assert(IsUuid(Uuid().String()), Equals, true)
}

func (s *MySuite) TestUuid(c *C) {
	result := Uuid().String()
	if err = validator.New().Var(result, "uuid"); err != nil {
		c.Error(err)
	}
	c.Log(result)
}

func (s *MySuite) TestCamel(c *C) {
	c.Assert(Camel("my_lab"), Equals, "MyLab")
	c.Assert(Camel("my-lab"), Equals, "MyLab")
	c.Assert(Camel("my lab"), Equals, "MyLab")
}

func (s *MySuite) TestSnake(c *C) {
	c.Assert(Snake("MyLab"), Equals, "my_lab")
	c.Assert(Snake("my-lab"), Equals, "my_lab")
	c.Assert(Snake("my lab"), Equals, "my_lab")
}

func (s *MySuite) TestKebab(c *C) {
	c.Assert(Kebab("MyLab"), Equals, "my-lab")
	c.Assert(Kebab("my_lab"), Equals, "my-lab")
	c.Assert(Kebab("my lab"), Equals, "my-lab")
}

func (s *MySuite) TestLimit(c *C) {
	result := Limit("The quick brown fox jumps over the lazy dog", 20)
	c.Assert(result, Equals, "The quick brown fox...")
}
