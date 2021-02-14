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
	Verify(value ...interface{}) bool
	Factory(value ...interface{})
	Destory(jti string, ack string) error
}

// Create authorization logic
//	@param `ctx` *gin.Context
//	@param `cookie` typ.Cookie
//	@param `claims` jwt.MapClaims
//	@param `refresh` RefreshTokenAPI refreshToken factory
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

// Verify authorization logic
//	@param `ctx` *gin.Context
//	@param `cookie` typ.Cookie
//	@param `refresh` RefreshTokenAPI refreshToken verification
func Verify(ctx *gin.Context, cookie typ.Cookie, refresh RefreshTokenAPI) (err error) {
	var value string
	if value, err = ctx.Cookie(cookie.Name); err != nil {
		return UserLoginError
	}
	var parseClaims jwt.MapClaims
	if parseClaims, err = tokenx.Verify(value, func(claims jwt.MapClaims) (jwt.MapClaims, error) {
		jti := claims["jti"].(string)
		ack := claims["ack"].(string)
		if result := refresh.Verify(jti, ack); !result {
			return nil, RefreshTokenExpired
		}
		for _, defaultClaim := range []string{"aud", "exp", "jti", "iat", "iss", "nbf", "sub"} {
			delete(claims, defaultClaim)
		}
		var token *tokenx.Token
		if token, err = tokenx.Make(claims, time.Hour); err != nil {
			return nil, err
		}
		cookie.Set(ctx, token.Value)
		return token.Claims, nil
	}); err != nil {
		return
	}
	ctx.Set("auth", parseClaims)
	return
}

// Authorization verification middleware
//	@param `cookie` typ.Cookie
//	@param `refresh` RefreshTokenAPI refreshToken verification
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

// Destroy authorization logic
//	@param `ctx` *gin.Context
//	@param `cookie` typ.Cookie
//	@param `refresh` RefreshTokenAPI refreshToken destory verification
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
//	@param `ctx` *gin.Context
func Get(ctx *gin.Context) (jwt.MapClaims, error) {
	val, exists := ctx.Get("auth")
	if !exists {
		return nil, UserLoginError
	}
	return val.(jwt.MapClaims), nil
}
