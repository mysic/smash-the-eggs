package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/xujiajun/nutsdb"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"smash-golden-eggs/service"
)

type loginForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
	//captcha string `form:"captcha" binding:"required"`
}

type passwordForm struct {
	NewPassword string `form:"new_password" binding:"required"`
	OriginPassword string `form:"origin_password" binding:"required"`
}

// Login 登入
func Login (c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token != "" {
		tokenString, _, err := service.Getting(token)
		if err == nil && tokenString.Valid {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": 0, "msg": "admin已经登入"})
			return
		}
	}
	var params loginForm
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":"参数错误",
			"data":err.Error(),
		})
		return
	}
	bucket := service.AdminBucket
	username := []byte(c.PostForm("username"))
	password := []byte(c.PostForm("password"))
	var dbPassword []byte
	conn := service.Conn
	if err := conn.View(func(tx *nutsdb.Tx) error {
		val, err := tx.Get(bucket, username)
		if err != nil {
			return err
		}
		dbPassword = val.Value
		return nil
	});err != nil {}
	err := bcrypt.CompareHashAndPassword(dbPassword, password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":"用户名或密码错误",
		})
		return 
	}
	token, err = service.Setting("admin")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK,gin.H{
			"code": -1,
			"msg":"登陆失败",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":"登录成功",
		"data":token,
	})

}

// Logout 登出
func Logout (c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":"成功登出",
	})
}

// Password 密码修改
func Password (c *gin.Context) {
	var params passwordForm
	var err error
	bucket := service.AdminBucket
	key := []byte("admin")
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":"参数错误",
			"data":err.Error(),
		})
		return
	}
	var dbPassword []byte
	newPassword := c.PostForm("new_password")
	originPassword := c.PostForm("origin_password")
	conn := service.Conn
	if err := conn.View(func(tx *nutsdb.Tx) error {
		val, err := tx.Get(bucket, key)
		if err != nil {
			return err
		}
		dbPassword = val.Value
		return nil
	}); err != nil {}
	err = bcrypt.CompareHashAndPassword(dbPassword, []byte(originPassword))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":"原密码错误",
			"data":err.Error(),
		})
		return
	}
	if err = conn.Update(func(tx *nutsdb.Tx) error {
		val, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		return tx.Put(bucket, key, val, 0)
	});err != nil {}

}
// Captcha todo 验证码
func Captcha (c *gin.Context) {

}
