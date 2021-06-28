package cookie

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Option struct {
	Name     string
	MaxAge   int           `yaml:"max_age"`
	Path     string        `yaml:"path"`
	Domain   string        `yaml:"domain"`
	Secure   bool          `yaml:"secure"`
	HttpOnly bool          `yaml:"http_only"`
	SameSite http.SameSite `yaml:"same_site"`
}

type Cookie struct {
	Option
}

func Initialize(option Option) *Cookie {
	return &Cookie{
		Option: option,
	}
}

func (x *Cookie) Set(c *gin.Context, name string, value string) {
	c.SetCookie(name, value, x.MaxAge, x.Path, x.Domain, x.Secure, x.HttpOnly)
	c.SetSameSite(x.SameSite)
}
