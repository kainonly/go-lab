package nats_accounts

import (
	"bytes"
	"github.com/nats-io/nats-server/v2/conf"
	"github.com/nats-io/nkeys"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"text/template"
)

func TestNkeys(t *testing.T) {
	user, _ := nkeys.CreateUser()
	data := []byte("Hello World")
	sig, _ := user.Sign(data)

	err := user.Verify(data, sig)
	assert.NoError(t, err)
	seed, _ := user.Seed()
	t.Log(string(seed))

	publicKey, _ := user.PublicKey()
	t.Log(publicKey)
}

type Account struct {
	Name  string
	Users []User
}

type User struct {
	NKey string
}

func TestNatsFactory(t *testing.T) {
	accounts := []Account{
		{
			Name: "weplanx",
			Users: []User{
				{NKey: "123456"},
				{NKey: "456789"},
			},
		},
		{
			Name: "example",
			Users: []User{
				{NKey: "asdasd"},
				{NKey: "dsadas"},
			},
		},
	}
	tmpl, err := template.ParseFiles("./account.tpl")
	assert.NoError(t, err)
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, accounts)
	assert.NoError(t, err)
	err = os.WriteFile("./accounts.conf", buf.Bytes(), os.ModePerm)
	assert.NoError(t, err)
}

func TestNatsParse(t *testing.T) {
	data, err := conf.ParseFile("./accounts.conf")
	assert.NoError(t, err)
	t.Log(data)
}
