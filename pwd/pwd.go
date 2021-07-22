package pwd

import (
	"github.com/alexedwards/argon2id"
)

// Make Use argon2id hash to generate user password
func Make(text string, options ...*argon2id.Params) (string, error) {
	option := argon2id.DefaultParams
	if len(options) != 0 {
		option = options[0]
	}
	return argon2id.CreateHash(text, option)
}

// Verify Verifying that a password matches a hash
func Verify(text string, hashText string) (bool, error) {
	return argon2id.ComparePasswordAndHash(text, hashText)
}
