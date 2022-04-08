package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"smash-golden-eggs/service"
)

func GameState() gin.HandlerFunc{
	return func(c *gin.Context) {
		if service.GameInstance.State != true {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code":-1,
				"msg":"游戏尚未开启！",
			})
			return
		}
		c.Next()
	}
}
