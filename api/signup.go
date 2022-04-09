package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"smash-golden-eggs/service"
)

//自定义一个字符串


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
	token := c.GetHeader("Authorization")
	// 判断是否已经登录
	if token != "" {
		tokenString, mobile, err := service.Getting(token)
		if err == nil && tokenString.Valid {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": 0, "msg": "手机号[" + mobile + "]已经登入"})
			return
		}
	}
	//token 1小时超时设置
	token, err := service.Setting(params.Mobile)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code":-1,
			"msg": "手机号[" + params.Mobile + "]登入失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":0,
		"msg": "手机号[" + params.Mobile + "]登入成功",
		"data":token,
	})
}
