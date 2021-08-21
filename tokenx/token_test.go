package tokenx

import (
	"github.com/dgrijalva/jwt-go"
	. "gopkg.in/check.v1"
	"testing"
	"time"
)

var err error

func Test(t *testing.T) { TestingT(t) }

type MySuite struct {
	token        *Token
	expiredToken *Token
}

var _ = Suite(&MySuite{})

func (s *MySuite) SetUpTest(c *C) {
	LoadKey([]byte(`hello`))
	claims := jwt.MapClaims{
		"iss":      "gin-extra",
		"aud":      []string{"tests"},
		"username": "kain",
	}
	if s.token, err = Make(claims, time.Hour*2); err != nil {
		c.Error(err)
	}
	if s.expiredToken, err = Make(claims, time.Second); err != nil {
		c.Error(err)
	}
}

func (s *MySuite) TestVerify(c *C) {
	var claims jwt.MapClaims
	if claims, err = Verify(s.token.Value, nil); err != nil {
		c.Error(err)
	}
	c.Assert(claims.Valid(), Equals, nil)
	c.Assert(claims["iss"], Equals, "gin-extra")
	c.Assert(claims["aud"], DeepEquals, []interface{}{"tests"})
	c.Assert(claims["username"], Equals, "kain")
	time.Sleep(time.Second * 2)
	if _, err = Verify(s.expiredToken.Value, nil); err != nil {
		c.Assert(err.(*jwt.ValidationError), NotNil)
	}
}
