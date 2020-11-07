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

func (c *mvc) AutoController(path string, controller interface{}) {
	typ := reflect.TypeOf(controller)
	val := reflect.ValueOf(controller)
	for i := 0; i < typ.NumMethod(); i++ {
		name := typ.Method(i).Name
		method := val.MethodByName(name).Interface()
		c.routes.POST(path+"/"+xstrings.FirstRuneToLower(name), c.Handle(method))
	}
}
