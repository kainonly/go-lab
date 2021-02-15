package rbacx

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-extra/authx"
	"strings"
)

type UserAPI interface {
	Get(username string) (result map[string]interface{})
}

type RoleAPI interface {
	Get(keys []string, mode string) *hashset.Set
}

type AclAPI interface {
	Get(key string, policy string) *hashset.Set
}

// Rbac verification middleware
//	@param `prefix` path prefix
//	@param `user` UserAPI
//	@param `role` RoleAPI
//	@param `acl` AclAPI
func Middleware(prefix string, user UserAPI, role RoleAPI, acl AclAPI) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		path := strings.Replace(ctx.Request.URL.Path, prefix, "", 1)
		acts := strings.Split(path, "/")
		var auth jwt.MapClaims
		if auth, err = authx.Get(ctx); err != nil {
			ctx.AbortWithStatusJSON(200, gin.H{
				"error": 1,
				"msg":   err.Error(),
			})
		}
		userData := user.Get(auth["user"].(string))
		roleKeys := userData["role"].([]interface{})
		keys := make([]string, len(roleKeys))
		for index, value := range roleKeys {
			keys[index] = value.(string)
		}
		roleAcl := role.Get(keys, "acl")
		if userData["acl"] != nil {
			roleAcl.Add(userData["acl"].([]interface{})...)
		}
		policyCursor := ""
		policyValues := []string{"0", "1"}
		for _, val := range policyValues {
			if roleAcl.Contains(acts[0] + ":" + val) {
				policyCursor = val
			}
		}
		if policyCursor == "" {
			ctx.AbortWithStatusJSON(200, gin.H{
				"error": 1,
				"msg":   "rbac invalid, policy is empty",
			})
		}
		scope := acl.Get(acts[0], policyCursor)
		if scope.Empty() {
			ctx.AbortWithStatusJSON(200, gin.H{
				"error": 1,
				"msg":   "rbac invalid, scope is empty",
			})
		}
		if !scope.Contains(acts[1]) {
			ctx.AbortWithStatusJSON(200, gin.H{
				"error": 1,
				"msg":   "rbac invalid, access denied",
			})
		}
		ctx.Next()
	}
}
