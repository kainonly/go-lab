package rbac

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-extra/authx"
	"strings"
)

type UserAPI interface {
	Get(ctx context.Context, username string) (result map[string]interface{})
}

type RoleAPI interface {
	Get(ctx context.Context, keys []string, mode string) *hashset.Set
}

type AclAPI interface {
	Get(ctx context.Context, key string, policy string) *hashset.Set
}

// Middleware rbac verification
func Middleware(prefix string, userAPI UserAPI, roleAPI RoleAPI, aclAPI AclAPI) gin.HandlerFunc {
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
		redisCtx := context.Background()
		user := auth["user"].(string)
		data := userAPI.Get(redisCtx, user)
		roles := data["role"].([]interface{})
		roleKeys := make([]string, len(roles))
		for index, value := range roles {
			roleKeys[index] = value.(string)
		}
		roleAcl := roleAPI.Get(redisCtx, roleKeys, "acl")
		if data["acl"] != nil {
			roleAcl.Add(data["acl"].([]interface{})...)
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
		scope := aclAPI.Get(redisCtx, acts[0], policyCursor)
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
