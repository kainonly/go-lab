package mvc

import (
	"github.com/gin-gonic/gin"
	"github.com/huandu/xstrings"
	"reflect"
)

type mvc struct {
	routes     *gin.RouterGroup
	dependency interface{}
}

func Factory(routes *gin.RouterGroup, dependency interface{}) *mvc {
	c := new(mvc)
	c.routes = routes
	c.dependency = dependency
	return c
}

func (c *mvc) Handle(handlersFn interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		switch method := handlersFn.(type) {
		case func() interface{}:
			ctx.JSON(200, method())
			break
		case func(ctx *gin.Context) interface{}:
			ctx.JSON(200, method(ctx))
			break
		case func(*gin.Context, interface{}) interface{}:
			ctx.JSON(200, method(ctx, c.dependency))
			break
		}
	}
}

type Auto struct {
	Path        string
	Controller  interface{}
	Middlewares []Middleware
}

type Middleware struct {
	Handle gin.HandlerFunc
	Only   []string
}

func (c *mvc) AutoController(auto Auto) {
	typ := reflect.TypeOf(auto.Controller)
	val := reflect.ValueOf(auto.Controller)
	for i := 0; i < typ.NumMethod(); i++ {
		name := typ.Method(i).Name
		method := val.MethodByName(name).Interface()
		var handlers []gin.HandlerFunc
		for _, middleware := range auto.Middlewares {
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
		handlers = append(handlers, c.Handle(method))
		c.routes.POST(auto.Path+"/"+xstrings.FirstRuneToLower(name), handlers...)
	}
}
