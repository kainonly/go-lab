package str

import (
	"github.com/google/uuid"
	"math/rand"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// generates a random string of the specified length
//	@param `length` string
//	@return string
func Random(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// generates a UUID (version 4)
//	@return uuid.UUID
func Uuid() uuid.UUID {
	return uuid.New()
}
