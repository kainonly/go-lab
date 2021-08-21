package str

import (
	"github.com/google/uuid"
	"github.com/huandu/xstrings"
	"math/rand"
	"time"
)

// Random generates a random string of the specified length
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

// Uuid generates a UUID (version 4)
func Uuid() uuid.UUID {
	return uuid.New()
}

// Camel converts the given string to CamelCase
func Camel(str string) string {
	return xstrings.ToCamelCase(str)
}

// Snake converts the given string to snake_case
func Snake(str string) string {
	return xstrings.ToSnakeCase(str)
}

// Kebab converts the given string to kebab-case
func Kebab(str string) string {
	return xstrings.ToKebabCase(str)
}

// Limit truncates the given string to the specified length
func Limit(str string, length int) string {
	return str[:length-1] + "..."
}
