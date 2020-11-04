package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"time"
)

type Option struct {
	Issuer   string
	Audience []string
	Expires  uint
}

type Handle func(option Option) (claims jwt.MapClaims, err error)

var (
	Key     []byte
	Options map[string]Option
	Method  jwt.SigningMethod = jwt.SigningMethodHS256
)

func Make(scene string, claims jwt.MapClaims) (tokenString string, err error) {
	option, exists := Options[scene]
	if !exists {
		err = fmt.Errorf("the [%v] scene does not exist", scene)
		return
	}
	claims["jti"] = uuid.New()
	claims["iat"] = time.Now().Unix()
	claims["iss"] = option.Issuer
	claims["aud"] = option.Audience
	claims["exp"] = time.Now().Add(time.Second * time.Duration(option.Expires)).Unix()
	token := jwt.NewWithClaims(Method, claims)
	tokenString, err = token.SignedString(Key)
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
