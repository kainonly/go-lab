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

func Verify(ctx *gin.Context) {

}

func AuthVerify(cookie typ.Cookie, refresh RefreshTokenAPI) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		var tokenString string
		if tokenString, err = ctx.Cookie(cookie.Name); err != nil {
			ctx.AbortWithStatusJSON(200, gin.H{
				"error": 1,
				"msg":   err.Error(),
			})
			return
		}
		if _, err := tokenx.Verify(tokenString, func(claims jwt.MapClaims) (jwt.MapClaims, error) {
			jti := claims["jti"].(string)
			ack := claims["ack"].(string)
			if result := refresh.Verify(jti, ack); !result {
				return nil, errors.New("refresh token verification expired")
			}
			var token *tokenx.Token
			if token, err = tokenx.Make(jwt.MapClaims{
				"jti": jti,
				"ack": ack,
			}, time.Hour); err != nil {
				return nil, err
			}
			cookie.Set(ctx, token.Value)
			return token.Claims, nil
		}); err != nil {
			ctx.AbortWithStatusJSON(200, gin.H{
				"error": 1,
				"msg":   err.Error(),
			})
			return
		}
		ctx.Next()
	}
}
