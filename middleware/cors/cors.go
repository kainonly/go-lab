package cors

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type Option struct {
	Origin        []string
	Method        []string
	AllowHeader   []string
	ExposedHeader []string
	MaxAge        int
	Credentials   bool
}

func Cors(option Option) gin.HandlerFunc {
	origin := strings.Join(option.Origin, ",")
	method := strings.Join(option.Method, ",")
	allowHeader := strings.Join(option.AllowHeader, ",")
	exposedHeader := strings.Join(option.ExposedHeader, ",")
	maxAge := strconv.Itoa(option.MaxAge)
	return func(ctx *gin.Context) {
		ctx.Header("access-control-allow-origin", origin)
		ctx.Header("Access-Control-Allow-Methods", method)
		ctx.Header("Access-Control-Allow-Headers", allowHeader)
		ctx.Header("Access-Control-Expose-Headers", exposedHeader)
		ctx.Header("Access-Control-Max-Age", maxAge)
		if option.Credentials {
			ctx.Header("Access-Control-Allow-Credentials", "true")
		}
		if ctx.Request.Method == "OPTIONS" {
			ctx.Status(200)
			return
		}
		ctx.Next()
	}
}
