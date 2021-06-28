package auth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-helper/cookie"
	"github.com/kainonly/gin-helper/str"
	"github.com/kainonly/gin-helper/tokenx"
	"time"
)

var (
	UserLoginError      = errors.New("please first authorize user login")
	RefreshTokenExpired = errors.New("refresh token verification expired")
)

type Option struct {
	Key       string   `yaml:"key"`
	Issuer    string   `yaml:"issuer"`
	Audience  []string `yaml:"audience"`
	NotBefore int64    `yaml:"not_before"`
	Expires   int64    `yaml:"expires"`
}

type auth struct {
	signKey    []byte
	signMethod jwt.SigningMethod
	iss        string
	aud        []string
	nbf        int64
	exp        time.Duration
	cookie     *cookie.Cookie
	refreshFn  RefreshFn
}

type Ext struct {
	Method    jwt.SigningMethod
	Cookie    *cookie.Cookie
	RefreshFn RefreshFn
}

type RefreshFn interface {
	Factory(values ...interface{})
	Verify(values ...interface{}) bool
	Renewal(values ...interface{})
	Destory(values ...interface{}) (err error)
}

func Factory(option Option, ext Ext) *auth {
	return &auth{
		signKey:    []byte(option.Key),
		signMethod: ext.Method,
		iss:        option.Issuer,
		aud:        option.Audience,
		nbf:        option.NotBefore,
		exp:        time.Duration(option.Expires) * time.Second,
		cookie:     ext.Cookie,
		refreshFn:  ext.RefreshFn,
	}
}

type CreateOptions func(*createOption)

type createOption struct {
	sub    string
	uid    interface{}
	data   interface{}
	c      *gin.Context
	cookie string
}

func SetPayload(sub string, uid interface{}, data map[string]interface{}) CreateOptions {
	return func(option *createOption) {
		option.sub = sub
		option.uid = uid
		option.data = data
	}
}

func UseCookie(c *gin.Context, cookie string) CreateOptions {
	return func(option *createOption) {
		option.c = c
		option.cookie = cookie
	}
}

// Create authorization logic
func (x *auth) Create(options ...CreateOptions) (tokenString string, err error) {
	option := createOption{}
	for _, apply := range options {
		apply(&option)
	}
	claims := jwt.MapClaims{
		"iat":  time.Now().Unix(),
		"nbf":  time.Now().Add(time.Second * time.Duration(x.nbf)).Unix(),
		"exp":  time.Now().Add(x.exp).Unix(),
		"jti":  str.Uuid().String(),
		"uid":  option.uid,
		"data": option.data,
	}
	token := jwt.NewWithClaims(x.signMethod, claims)
	if tokenString, err = token.SignedString(x.signKey); err != nil {
		return
	}
	if option.c != nil && option.cookie != "" {
		x.cookie.Set(option.c, option.cookie, tokenString)
	}
	if x.refreshFn != nil {
		x.refreshFn.Factory(claims)
	}
	return
}

// Verify authorization logic
func (x *auth) Verify(tokenString string) (err error) {
	var token *jwt.Token
	if token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return x.signKey, nil
	}); err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors == jwt.ValidationErrorExpired && x.refreshFn != nil && token != nil {
				claims := token.Claims.(jwt.MapClaims)
				if result := x.refreshFn.Verify(claims); !result {
					return RefreshTokenExpired
				}
			}
		}
		return
	}
	if x.refreshFn != nil {
		x.refreshFn.Renewal(token.Claims)
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
