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
	acls  AclFn
	roles RoleFn
	users UserFn
}

type AclFn interface {
	Fetch(ctx context.Context, key string, mode string) (*hashset.Set, error)
}

type RoleMode string

const (
	RoleAcl        RoleMode = "acl"
	RoleResource   RoleMode = "resource"
	RolePermission RoleMode = "permission"
)

type RoleFn interface {
	Fetch(ctx context.Context, keys []string, mode RoleMode) (*hashset.Set, error)
}

type UserFn interface {
	Fetch(ctx context.Context, uid interface{}) (map[string]interface{}, error)
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
			return
		}
		mClaims := claims.(jwt.MapClaims)
		scope := Scope{
			UID:  mClaims["uid"],
			Data: mClaims["data"].(map[string]interface{}),
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		var err error
		var user map[string]interface{}
		if user, err = fn.users.Fetch(ctx, scope.UID); err != nil {
			c.AbortWithStatusJSON(200, gin.H{
				"error": 1,
				"msg":   err.Error(),
			})
			return
		}
		roles := user["role"].([]interface{})
		roleKeys := make([]string, len(roles))
		for index, value := range roles {
			roleKeys[index] = value.(string)
		}
		var roleAcl *hashset.Set
		if roleAcl, err = fn.roles.Fetch(ctx, roleKeys, RoleAcl); err != nil {
			c.AbortWithStatusJSON(200, gin.H{
				"error": 1,
				"msg":   err.Error(),
			})
			return
		}
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
			c.AbortWithStatusJSON(200, gin.H{
				"error": 1,
				"msg":   "rbac invalid, policy is empty",
			})
			return
		}
		var acl *hashset.Set
		if acl, err = fn.acls.Fetch(ctx, acts[0], policyCursor); err != nil {
			c.AbortWithStatusJSON(200, gin.H{
				"error": 1,
				"msg":   err.Error(),
			})
			return
		}
		if acl.Empty() {
			c.AbortWithStatusJSON(200, gin.H{
				"error": 1,
				"msg":   "rbac invalid, scope is empty",
			})
			return
		}
		if !acl.Contains(acts[1]) {
			c.AbortWithStatusJSON(200, gin.H{
				"error": 1,
				"msg":   "rbac invalid, access denied",
			})
			return
		}
		c.Next()
	}
}
