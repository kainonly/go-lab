package mvcx

import (
	"github.com/gin-gonic/gin"
	"github.com/huandu/xstrings"
	"reflect"
)

type mvcx struct {
	routes     *gin.RouterGroup
	dependency interface{}
}

func Initialize(routes *gin.RouterGroup, dependency interface{}) *mvcx {
	return &mvcx{
		routes:     routes,
		dependency: dependency,
	}
}

type Middleware struct {
	Handle gin.HandlerFunc
	Only   []string
}

func Handle(handlerFn interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if fn, ok := handlerFn.(func(ctx *gin.Context) interface{}); ok {
			switch result := fn(ctx).(type) {
			case bool:
				if result {
					ctx.JSON(200, gin.H{
						"error": 0,
						"msg":   "ok",
					})
				} else {
					ctx.JSON(200, gin.H{
						"error": 1,
						"msg":   "failed",
					})
				}
				break
			case error:
				ctx.JSON(200, gin.H{
					"error": 1,
					"msg":   result.Error(),
				})
				break
			default:
				ctx.JSON(200, gin.H{
					"error": 0,
					"data":  result,
				})
			}
		} else {
			ctx.Status(404)
		}
	}
}

func (c *mvcx) AutoController(path string, controller interface{}, middlewares ...Middleware) {
	if control, ok := controller.(interface {
		Inject(dependency interface{})
	}); ok {
		control.Inject(c.dependency)
	}
	typ := reflect.TypeOf(controller)
	val := reflect.ValueOf(controller)
	for i := 0; i < typ.NumMethod(); i++ {
		name := typ.Method(i).Name
		method := val.MethodByName(name).Interface()
		var handlers []gin.HandlerFunc
		for _, middleware := range middlewares {
			if len(middleware.Only) == 0 {
				handlers = append(handlers, middleware.Handle)
			} else {
				for _, m := range middleware.Only {
					if m == name {
						handlers = append(handlers, middleware.Handle)
					}
				}
			}
		}
		handlers = append(handlers, Handle(method))
		c.routes.POST(path+"/"+xstrings.FirstRuneToLower(name), handlers...)
	}
}