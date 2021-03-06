package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"smash-golden-eggs/service"
)

func OrderState() gin.HandlerFunc{
	return func(c *gin.Context) {
		if service.OrderState != service.OrderStatePaid {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code":-1,
				"msg":"订单还未支付",
			})
			return
		}
		c.Next()
	}
}