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

var (
	UserLoginError      = errors.New("please first authorize user login")
	RefreshTokenExpired = errors.New("refresh token verification expired")
)

type RefreshTokenAPI interface {
	Factory(value ...interface{})
	Renewal(value ...interface{})
	Verify(value ...interface{}) bool
	Destory(value ...interface{}) (err error)
}

// Create authorization logic
func Create(ctx *gin.Context, cookie typ.Cookie, claims jwt.MapClaims, refresh RefreshTokenAPI) (err error) {
	claims["jit"] = str.Uuid().String()
	claims["ack"] = str.Random(8)
	var token *tokenx.Token
	if token, err = tokenx.Make(claims, time.Hour); err != nil {
		return
	}
	refresh.Factory(claims["jit"], claims["ack"], time.Hour)
	cookie.Set(ctx, token.Value)
	return
}

// Verify authorization logic
func Verify(ctx *gin.Context, cookie typ.Cookie, refresh RefreshTokenAPI) (err error) {
	var value string
	if value, err = ctx.Cookie(cookie.Name); err != nil {
		return UserLoginError
	}
	var parseClaims jwt.MapClaims
	if parseClaims, err = tokenx.Verify(value, func(claims jwt.MapClaims) (jwt.MapClaims, error) {
		if result := refresh.Verify(claims["jti"].(string), claims["ack"].(string)); !result {
			return nil, RefreshTokenExpired
		}
		for _, defaultClaim := range []string{"aud", "exp", "jti", "iat", "iss", "nbf", "sub"} {
			delete(claims, defaultClaim)
		}
		var token *tokenx.Token
		if token, err = tokenx.Make(claims, time.Minute*15); err != nil {
			return nil, err
		}
		cookie.Set(ctx, token.Value)
		return token.Claims, nil
	}); err != nil {
		return
	}
	refresh.Renewal(parseClaims["jti"].(string), time.Hour)
	ctx.Set("auth", parseClaims)
	return
}

// Middleware authorization verification
func Middleware(cookie typ.Cookie, refresh RefreshTokenAPI) gin.HandlerFunc {
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

// Destory authorization logic
func Destory(ctx *gin.Context, cookie string, refresh RefreshTokenAPI) (err error) {
	var value string
	if value, err = ctx.Cookie(cookie); err != nil {
		return nil
	}
	var claims jwt.MapClaims
	if claims, err = tokenx.Verify(value, func(c jwt.MapClaims) (jwt.MapClaims, error) {
		return c, nil
	}); err != nil {
		return
	}
	return refresh.Destory(claims["jti"].(string), claims["ack"].(string))
}

// Get authorization claims
func Get(ctx *gin.Context) (jwt.MapClaims, error) {
	val, exists := ctx.Get("auth")
	if !exists {
		return nil, UserLoginError
	}
	return val.(jwt.MapClaims), nil
}
