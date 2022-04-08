package api

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type mobileForm struct {
	Mobile string `form:"mobile" binding:"required"`
}

func MobileSignUp(c *gin.Context) {
	var params mobileForm
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":"参数错误",
			"data":err.Error(),
		})
		return
	}
	session := sessions.Default(c)
	// 判断是否已经登录
	mobile := session.Get("mobile")
	if mobile != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":0,
			"msg": "手机号[" + mobile.(string) + "]已经登入",
		})
		return
	}
	//session 1小时超时设置
	session.Options(sessions.Options{MaxAge: 3600})
	session.Set("mobile",params.Mobile)
	err := session.Save()
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":-1,
			"msg": "手机号[" + params.Mobile + "]登入失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":0,
		"msg": "手机号[" + params.Mobile + "]登入成功",
	})
}
