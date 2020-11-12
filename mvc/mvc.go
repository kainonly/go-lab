package mvc

import (
	"github.com/gin-gonic/gin"
	"github.com/huandu/xstrings"
	"reflect"
)

type mvc struct {
	routes *gin.RouterGroup
}

func Factory(routes *gin.RouterGroup) *mvc {
	c := new(mvc)
	c.routes = routes
	return c
}

func Handle(handlersFn interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if method, ok := handlersFn.(func(ctx *gin.Context) interface{}); ok {
			result := method(ctx)
			switch val := result.(type) {
			case bool:
				if val {
					ctx.Status(200)
				} else {
					ctx.Status(500)
				}
				break
			case error:
				ctx.JSON(400, val.Error())
				break
			default:
				ctx.JSON(200, val)
			}
		} else {
			ctx.Status(404)
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
		handlers = append(handlers, Handle(method))
		c.routes.POST(auto.Path+"/"+xstrings.FirstRuneToLower(name), handlers...)
	}
}
