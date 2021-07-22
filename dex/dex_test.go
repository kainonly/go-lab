package dex

import (
	"log"
	"os"
	"testing"
)

var x *Dex
var err error

func TestMain(m *testing.M) {
	x, err = Make(Option{Key: "6ixSiEXaqxsJTozbnxQ76CWdZXB2JazK"})
	if err != nil {
		log.Fatalln(err)
	}
	os.Exit(m.Run())
}

func TestDexId(t *testing.T) {
	hash, err := x.EncodeId([]int{651})
	if err != nil {
		t.Error(err)
	}
	t.Log(hash)
	val, err := x.DecodeId(hash)
	if err != nil {
		t.Error(err)
	}
	t.Log(val)
}

func TestDexData(t *testing.T) {
	data := []byte("Gophers, gophers, gophers everywhere!")
	ciphertext, err := x.Encode(data)
	if err != nil {
		t.Error(err)
	}
	t.Log(ciphertext)
	result, err := x.Decode(ciphertext)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(result))
}
