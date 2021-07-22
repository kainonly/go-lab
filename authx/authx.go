package authx

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-helper/cookie"
	"github.com/kainonly/gin-helper/str"
	"time"
)

var (
	UserLoginError      = errors.New("当前鉴权已失效")
	RefreshTokenExpired = errors.New("刷新令牌已过期")
)

type Option struct {
	Key       string   `yaml:"key"`
	Issuer    string   `yaml:"issuer"`
	Audience  []string `yaml:"audience"`
	NotBefore int64    `yaml:"not_before"`
	Expires   int64    `yaml:"expires"`
}

type Auth struct {
	signKey    []byte
	signMethod jwt.SigningMethod
	iss        string
	aud        []string
	nbf        int64
	exp        time.Duration
	cookie     *cookie.Cookie
	refreshFn  RefreshFn
}

type Args struct {
	Method    jwt.SigningMethod
	UseCookie *cookie.Cookie
	RefreshFn RefreshFn
}

type RefreshFn interface {
	Factory(claims jwt.MapClaims, args ...interface{})
	Verify(claims jwt.MapClaims, args ...interface{}) bool
	Renewal(claims jwt.MapClaims, args ...interface{})
	Destory(claims jwt.MapClaims, args ...interface{}) (err error)
}

func Make(option Option, args Args) *Auth {
	return &Auth{
		signKey:    []byte(option.Key),
		signMethod: args.Method,
		iss:        option.Issuer,
		aud:        option.Audience,
		nbf:        option.NotBefore,
		exp:        time.Duration(option.Expires) * time.Second,
		cookie:     args.UseCookie,
		refreshFn:  args.RefreshFn,
	}
}

// Create 创建认证
func (x *Auth) Create(c *gin.Context, sub interface{}, uid interface{}, data interface{}) (raw string, err error) {
	claims := jwt.MapClaims{
		"iat":  time.Now().Unix(),
		"nbf":  time.Now().Add(time.Second * time.Duration(x.nbf)).Unix(),
		"exp":  time.Now().Add(x.exp).Unix(),
		"jti":  str.Uuid().String(),
		"sub":  sub,
		"uid":  uid,
		"data": data,
	}
	token := jwt.NewWithClaims(x.signMethod, claims)
	if raw, err = token.SignedString(x.signKey); err != nil {
		return
	}
	if x.cookie != nil {
		x.cookie.Set(c, raw)
	}
	if x.refreshFn != nil {
		x.refreshFn.Factory(claims)
	}
	c.Set("claims", claims)
	return
}

// Verify 鉴权认证
func (x *Auth) Verify(c *gin.Context, args ...interface{}) (err error) {
	var raw string
	if x.cookie != nil {
		if raw, err = c.Cookie(x.cookie.Name); err != nil {
			return UserLoginError
		}
	} else {
		if len(args) != 0 {
			raw = args[0].(string)
		}
	}
	if raw == "" {
		return UserLoginError
	}
	var token *jwt.Token
	if token, err = jwt.Parse(raw, func(token *jwt.Token) (interface{}, error) {
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
				updateClaims := jwt.MapClaims{
					"iat":  time.Now().Unix(),
					"nbf":  time.Now().Add(time.Second * time.Duration(x.nbf)).Unix(),
					"exp":  time.Now().Add(x.exp).Unix(),
					"jti":  str.Uuid().String(),
					"sub":  claims["sub"],
					"uid":  claims["uid"],
					"data": claims["data"],
				}
				token = jwt.NewWithClaims(x.signMethod, updateClaims)
				if raw, err = token.SignedString(x.signKey); err != nil {
					return
				}
				if x.cookie != nil {
					x.cookie.Set(c, raw)
				}
				c.Set("token", token)
			}
		}
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	if x.refreshFn != nil {
		x.refreshFn.Renewal(claims)
	}
	c.Set("claims", claims)
	return
}

// Destory 销毁认证
func (x *Auth) Destory(c *gin.Context, args ...interface{}) (err error) {
	if err = x.Verify(c, args); err != nil {
		return
	}
	claims, exists := c.Get("claims")
	if !exists {
		return fmt.Errorf("environment verification is abnormal")
	}
	if x.cookie != nil {
		x.cookie.Set(c, "")
	}
	if x.refreshFn != nil {
		if err = x.refreshFn.Destory(claims.(jwt.MapClaims)); err != nil {
			return
		}
	}
	return
}

// Middleware 鉴权认证中间件
func Middleware(auth Auth, args ...interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := auth.Verify(c, args); err != nil {
			c.AbortWithStatusJSON(200, gin.H{
				"error": 0,
				"msg":   err.Error(),
			})
			return
		}
		c.Next()
	}
}
