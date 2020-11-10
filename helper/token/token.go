package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"time"
)

var (
	Key     []byte
	Options map[string]Option
	Method  jwt.SigningMethod = jwt.SigningMethodHS256
)

type Token struct {
	value  string
	claims jwt.MapClaims
}

func (c *Token) String() string {
	return c.value
}

func (c *Token) Claims() jwt.MapClaims {
	return c.claims
}

type Option struct {
	Issuer   string   `yaml:"issuer"`
	Audience []string `yaml:"audience"`
	Expires  uint     `yaml:"expires"`
}

type Handle func(option Option) (claims jwt.MapClaims, err error)

func Make(scene string, claims jwt.MapClaims) (token *Token, err error) {
	option, exists := Options[scene]
	if !exists {
		err = fmt.Errorf("the [%v] scene does not exist", scene)
		return
	}
	token = new(Token)
	claims["jti"] = uuid.New()
	claims["iat"] = time.Now().Unix()
	claims["iss"] = option.Issuer
	claims["aud"] = option.Audience
	claims["exp"] = time.Now().Add(time.Second * time.Duration(option.Expires)).Unix()
	t := jwt.NewWithClaims(Method, claims)
	token.claims = t.Claims.(jwt.MapClaims)
	token.value, err = t.SignedString(Key)
	if err != nil {
		return
	}

	return
}

func Verify(scene string, tokenString string, refresh Handle) (claims jwt.MapClaims, err error) {
	option, exists := Options[scene]
	if !exists {
		err = fmt.Errorf("the [%v] scene does not exist", scene)
		return
	}
	var token *jwt.Token
	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return Key, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors == jwt.ValidationErrorExpired {
				return refresh(option)
			}
		}
	} else {
		claims = token.Claims.(jwt.MapClaims)
	}
	return
}
