package token

import (
	"github.com/dgrijalva/jwt-go"
	"testing"
	"time"
)

var token *Token
var err error

func TestMake(t *testing.T) {
	claims := jwt.MapClaims{
		"iss":      "gin-extra",
		"aud":      []string{"tests"},
		"username": "kain",
	}
	if token, err = Make(claims, time.Hour*2); err != nil {
		t.Fatal(err)
	}
	t.Log(token.String())
}

func TestVerify(t *testing.T) {
	claims, err := Verify(token.String(), nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(claims)
}
