package mvc

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Ok gin.H
type Create gin.H

func Bind(handlerFn interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if fn, ok := handlerFn.(func(ctx *gin.Context) interface{}); ok {
			switch result := fn(ctx).(type) {
			case Ok:
				ctx.JSON(http.StatusOK, result)
				break
			case Create:
				ctx.JSON(http.StatusCreated, result)
				break
			case error:
				ctx.JSON(http.StatusBadRequest, gin.H{
					"msg": result.Error(),
				})
				break
			default:
				ctx.JSON(http.StatusNoContent, gin.H{
					"msg": "ok",
				})
			}
		}
	}
}
