package token

import (
	"github.com/dgrijalva/jwt-go"
	"os"
	"testing"
)

var token *Token
var err error

func TestMain(m *testing.M) {
	Options = map[string]Option{
		"system": {
			Issuer:   "helper",
			Audience: []string{"dev"},
			Expires:  3600,
		},
	}
	os.Exit(m.Run())
}

func TestMake(t *testing.T) {
	claims := jwt.MapClaims{
		"username": "kain",
	}
	if token, err = Make("system", claims); err != nil {
		t.Fatal(err)
	}
	t.Log(token.String())
}

func TestVerify(t *testing.T) {
	claims, err := Verify("system", token.String(), nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(claims)
}
