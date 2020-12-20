package authx

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-extra/str"
	"github.com/kainonly/gin-extra/tokenx"
	"github.com/kainonly/gin-extra/typ"
	"time"
)

type RefreshTokenAPI interface {
	Verify(value ...interface{}) bool
	Factory(value ...interface{})
}

func Create(ctx *gin.Context, cookie typ.Cookie, claims jwt.MapClaims, refresh RefreshTokenAPI) (err error) {
	jti := str.Uuid()
	ack := str.Random(8)
	defaultClaims := jwt.MapClaims{
		"jti": jti,
		"ack": ack,
	}
	for key, value := range claims {
		defaultClaims[key] = value
	}
	var token *tokenx.Token
	if token, err = tokenx.Make(defaultClaims, time.Hour); err != nil {
		return
	}
	refresh.Factory(jti.String(), ack, time.Hour*72)
	cookie.Set(ctx, token.Value)
	return
}

func Verify(ctx *gin.Context, cookie typ.Cookie, refresh RefreshTokenAPI) (err error) {
	var value string
	if value, err = ctx.Cookie(cookie.Name); err != nil {
		return
	}
	if _, err = tokenx.Verify(value, func(claims jwt.MapClaims) (jwt.MapClaims, error) {
		jti := claims["jti"].(string)
		ack := claims["ack"].(string)
		if result := refresh.Verify(jti, ack); !result {
			return nil, errors.New("refresh token verification expired")
		}
		defaultClaims := jwt.MapClaims{
			"jti": jti,
			"ack": ack,
		}
		standardClaims := []string{"aud", "exp", "jti", "iat", "iss", "nbf", "sub"}
		for key, value := range claims {
			for _, claimName := range standardClaims {
				if key == claimName {
					continue
				}
			}
			defaultClaims[key] = value
		}
		var token *tokenx.Token
		if token, err = tokenx.Make(defaultClaims, time.Hour); err != nil {
			return nil, err
		}
		cookie.Set(ctx, token.Value)
		return token.Claims, nil
	}); err != nil {
		return
	}
	return
}

func AuthVerify(cookie typ.Cookie, refresh RefreshTokenAPI) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := Verify(ctx, cookie, refresh); err != nil {
			ctx.AbortWithStatusJSON(200, gin.H{
				"error": 1,
				"msg":   err.Error(),
			})
			return
		}
		ctx.Next()
	}
}
