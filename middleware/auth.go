package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		var isAdmin interface{}
		if isAdmin = session.Get("mobile"); isAdmin == nil {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code":-1,
				"msg":"未登入，请先登入",
			})
			return
		}
		c.Next()
	}
}

func AdminAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		var isAdmin interface{}
		isAdmin = session.Get("isAdmin")
		if isAdmin == nil {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code":-1,
				"msg":"未登入，请先登入",
			})
			return
		}

		if isAdmin.(string) != "yes" {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code":-1,
				"msg":"不是管理员，没有操作权限",
			})
			return
		}
		c.Next()
	}
}

