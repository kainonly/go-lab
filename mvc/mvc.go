package mvc

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Ok gin.H
type Create gin.H

func Bind(handlerFn interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if fn, ok := handlerFn.(func(c *gin.Context) interface{}); ok {
			switch result := fn(c).(type) {
			case Ok:
				c.JSON(http.StatusOK, result)
				break
			case Create:
				c.JSON(http.StatusCreated, result)
				break
			case error:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": result.Error(),
				})
				break
			default:
				c.JSON(http.StatusNoContent, gin.H{
					"msg": "ok",
				})
			}
		}
	}
}
