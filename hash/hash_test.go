package hash

import (
	"github.com/alexedwards/argon2id"
	"testing"
)

var checkHash string

func TestMake(t *testing.T) {
	hash, err := Make(`pass`)
	if err != nil {
		t.Error(err)
	}
	t.Log(hash)
	checkHash = hash
}

func TestCheck(t *testing.T) {
	result, err := Verify(
		`pass`,
		checkHash,
	)
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}

func TestMakeExtend(t *testing.T) {
	option := argon2id.DefaultParams
	option.Memory = 128 * 1024
	option.Iterations = 6
	option.Parallelism = 2
	hash, err := Make(`pass`, option)
	if err != nil {
		t.Error(err)
	}
	t.Log(hash)
	checkHash = hash
}

func TestCheckExtend(t *testing.T) {
	result, err := Verify(
		`pass`,
		checkHash,
	)
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}
