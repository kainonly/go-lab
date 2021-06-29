package cookie

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Cookie struct {
	Option
}

type Option struct {
	Name     string
	MaxAge   int    `yaml:"max_age"`
	Path     string `yaml:"path"`
	Domain   string `yaml:"domain"`
	Secure   bool   `yaml:"secure"`
	HttpOnly bool   `yaml:"http_only"`
	SameSite http.SameSite
}

func Make(option Option, samesite http.SameSite) *Cookie {
	option.SameSite = samesite
	return &Cookie{
		Option: option,
	}
}

func (x *Cookie) Set(c *gin.Context, value string) {
	c.SetCookie(x.Name, value, x.MaxAge, x.Path, x.Domain, x.Secure, x.HttpOnly)
	c.SetSameSite(x.SameSite)
}
