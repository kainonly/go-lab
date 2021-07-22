package dex

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"github.com/speps/go-hashids/v2"
	"golang.org/x/crypto/chacha20poly1305"
)

type Option struct {
	Key string `yaml:"key"`
}

type Dex struct {
	aead   cipher.AEAD
	hashId *hashids.HashID
}

// Make 初始化数据加密
func Make(option Option) (x *Dex, err error) {
	x = new(Dex)
	if x.aead, err = chacha20poly1305.NewX([]byte(option.Key)); err != nil {
		return
	}
	hd := hashids.NewData()
	hd.Salt = option.Key
	if x.hashId, err = hashids.NewWithData(hd); err != nil {
		return
	}
	return
}

// EncodeId ID 加密
func (x *Dex) EncodeId(value []int) (string, error) {
	return x.hashId.Encode(value)
}

// DecodeId ID 解密
func (x *Dex) DecodeId(value string) ([]int, error) {
	return x.hashId.DecodeWithError(value)
}

// Encode 数据加密
func (x *Dex) Encode(data []byte) (string, error) {
	nonce := make([]byte, x.aead.NonceSize(), x.aead.NonceSize()+len(data)+x.aead.Overhead())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	encrypted := x.aead.Seal(nonce, nonce, data, nil)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// Decode 数据解密
func (x *Dex) Decode(text string) ([]byte, error) {
	encrypted, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return nil, err
	}
	nonce, ciphertext := encrypted[:x.aead.NonceSize()], encrypted[x.aead.NonceSize():]
	return x.aead.Open(nil, nonce, ciphertext, nil)
}
