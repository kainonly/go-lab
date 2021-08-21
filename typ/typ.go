package typ

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Cookie struct {
	Name     string        `yaml:"name"`
	MaxAge   int           `yaml:"max_age"`
	Path     string        `yaml:"path"`
	Domain   string        `yaml:"domain"`
	Secure   bool          `yaml:"secure"`
	HttpOnly bool          `yaml:"http_only"`
	SameSite http.SameSite `yaml:"same_site"`
}

func (c *Cookie) Set(ctx *gin.Context, value string) {
	ctx.SetCookie(c.Name, value, c.MaxAge, c.Path, c.Domain, c.Secure, c.HttpOnly)
	ctx.SetSameSite(c.SameSite)
}
