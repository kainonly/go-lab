package hash

import (
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

func TestMake2(t *testing.T) {
	hash, err := Make(`pass`, Time(6), Memory(128*1024), Threads(2))
	if err != nil {
		t.Error(err)
	}
	t.Log(hash)
	checkHash = hash
}

func TestCheck2(t *testing.T) {
	result, err := Verify(
		`pass`,
		checkHash,
	)
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}
