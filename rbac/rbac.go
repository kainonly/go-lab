package rbac

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type Fn struct {
	users UserFn
	roles RoleFn
	acls  AclFn
}

type UserFn interface {
	Get(ctx context.Context, uid interface{}) map[string]interface{}
}

type RoleMode int

const (
	RoleAcl      RoleMode = 0
	RoleResource RoleMode = 1
)

type RoleFn interface {
	Get(ctx context.Context, keys []string, mode RoleMode) *hashset.Set
}

type AclFn interface {
	Get(ctx context.Context, key string, policy string) *hashset.Set
}

type Scope struct {
	UID  interface{}
	Data map[string]interface{}
}

// Middleware rbac verification
func Middleware(prefix string, fn Fn) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := strings.Replace(c.Request.URL.Path, prefix, "", 1)
		acts := strings.Split(path, "/")
		claims, exists := c.Get("claims")
		if !exists {
			c.AbortWithStatusJSON(400, gin.H{
				"msg": "environment verification is abnormal",
			})
		}
		mClaims := claims.(jwt.MapClaims)
		scope := Scope{
			UID:  mClaims["uid"],
			Data: mClaims["data"].(map[string]interface{}),
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		user := fn.users.Get(ctx, scope.UID)
		roles := user["role"].([]interface{})
		roleKeys := make([]string, len(roles))
		for index, value := range roles {
			roleKeys[index] = value.(string)
		}
		roleAcl := fn.roles.Get(ctx, roleKeys, RoleAcl)
		if user["acl"] != nil {
			roleAcl.Add(user["acl"].([]interface{})...)
		}
		policyCursor := ""
		policyValues := []string{"0", "1"}
		for _, val := range policyValues {
			if roleAcl.Contains(acts[0] + ":" + val) {
				policyCursor = val
			}
		}
		if policyCursor == "" {
			c.AbortWithStatusJSON(400, gin.H{
				"msg": "rbac invalid, policy is empty",
			})
		}
		acl := fn.acls.Get(ctx, acts[0], policyCursor)
		if acl.Empty() {
			c.AbortWithStatusJSON(400, gin.H{
				"msg": "rbac invalid, scope is empty",
			})
		}
		if !acl.Contains(acts[1]) {
			c.AbortWithStatusJSON(400, gin.H{
				"msg": "rbac invalid, access denied",
			})
		}
		c.Next()
	}
}
