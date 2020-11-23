package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"time"
)

var (
	Key    []byte
	Method jwt.SigningMethod = jwt.SigningMethodHS256
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

// create a token
// @param claims jwt.MapClaims @refer http://self-issued.info/docs/draft-ietf-oauth-json-web-token.html#Claims
// @param expires time.Duration
// @return token *Token
func Make(claims jwt.MapClaims, expires time.Duration) (token *Token, err error) {
	token = new(Token)
	claims["jti"] = uuid.New()
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(expires).Unix()
	ref := jwt.NewWithClaims(Method, claims)
	token.claims = ref.Claims.(jwt.MapClaims)
	if token.value, err = ref.SignedString(Key); err != nil {
		return
	}
	return
}

// token refresh logic
// @param claims jwt.MapClaims
// @return jwt.MapClaims
type RefreshHandle func(claims jwt.MapClaims) (jwt.MapClaims, error)

// verify that the token is valid
// @param tokenString string
// @param refresh RefreshHandle
// @return claims jwt.MapClaims
func Verify(tokenString string, refresh RefreshHandle) (claims jwt.MapClaims, err error) {
	var token *jwt.Token
	if token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return Key, nil
	}); err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors == jwt.ValidationErrorExpired {
				if token != nil {
					return refresh(token.Claims.(jwt.MapClaims))
				}
			}
		}
		return
	}
	return token.Claims.(jwt.MapClaims), nil
}
