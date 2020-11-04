package token

import (
	"github.com/dgrijalva/jwt-go"
	"os"
	"testing"
)

var tokenString string
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
	tokenString, err = Make("system", jwt.MapClaims{
		"username": "kain",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tokenString)
}

func TestVerify(t *testing.T) {
	claims, err := Verify("system", tokenString, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(claims)
}
