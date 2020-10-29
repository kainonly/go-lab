package token

import (
	"os"
	"testing"
)

var token []byte
var err error

func TestMain(m *testing.M) {
	Options = map[string]Option{
		"system": {
			Issuer:   "iris-helper",
			Audience: []string{"tester"},
			Expires:  3600,
		},
	}
	os.Exit(m.Run())
}

func TestMake(t *testing.T) {
	token, err = Make("system", map[string]interface{}{
		"username": "kain",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(token))
}

func TestVerify(t *testing.T) {
	claims, err := Verify("system", token, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(claims)
}
