package str

import (
	"github.com/google/uuid"
	"github.com/huandu/xstrings"
	"math/rand"
	"time"
)

// generates a random string of the specified length
//	@param `length` string
//	@return string
func Random(length int, letterRunes ...rune) string {
	b := make([]rune, length)
	if len(letterRunes) == 0 {
		letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	}
	rand.Seed(time.Now().UnixNano())
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

// converts the given string to CamelCase
//	@param `str` string
//	@return string
func Camel(str string) string {
	return xstrings.ToCamelCase(str)
}

// converts the given string to snake_case
//	@param `str` string
//	@return string
func Snake(str string) string {
	return xstrings.ToSnakeCase(str)
}

// converts the given string to kebab-case
//	@param `str` string
//	@return string
func Kebab(str string) string {
	return xstrings.ToKebabCase(str)
}

// truncates the given string to the specified length
//	@param `str` string
//	@param `length` int
//	@return string
func Limit(str string, length int) string {
	return str[:length-1] + "..."
}
