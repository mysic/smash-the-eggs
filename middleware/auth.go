package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"smash-golden-eggs/service"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": -1, "msg": "未登入，请先登入"})
			return
		}
		tokenString, _, err := service.Getting(token)
		if err != nil || !tokenString.Valid {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": -1, "msg": "未登入，请先登入"})
			return
		}
		c.Next()
	}
}

