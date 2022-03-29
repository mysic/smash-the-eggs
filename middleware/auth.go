package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"smash-golden-eggs/service"
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
		if service.AdminState == false {

			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code":-1,
				"msg":"未登入，请先登入",
			})
			return
		}

		c.Next()
	}
}

