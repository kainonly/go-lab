package cors

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type Option struct {

	// Matches the request origin
	Origin []string `yaml:"origin"`

	// Matches the request method
	Method []string `yaml:"method"`

	// Sets the Access-Control-Allow-Headers response header
	AllowHeader []string `yaml:"allow_header"`

	// Sets the Access-Control-Expose-Headers response header
	ExposedHeader []string `yaml:"exposed_header"`

	// Sets the Access-Control-Max-Age response header
	MaxAge int `yaml:"max_age"`

	// Sets the Access-Control-Allow-Credentials header
	Credentials bool `yaml:"credentials"`
}

// Adds CORS (Cross-Origin Resource Sharing) headers support in your Gin application
//	@param `option` Option
//	@return gin.HandlerFunc
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
