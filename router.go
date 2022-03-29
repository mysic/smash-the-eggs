package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xujiajun/nutsdb"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"smash-golden-eggs/admin"
	"smash-golden-eggs/api"
	"smash-golden-eggs/middleware"
	"smash-golden-eggs/service"
	"time"
)

func RegisterRouter(r *gin.Engine) {

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":0,
			"msg":"welcome",
		})
	})

	apiRouter := r.Group("/api")
	{
		//注册手机号，建立会话

		apiRouter.POST("/signup", api.MobileSignUp)
		//获取本轮游戏内容
		apiRouter.GET("/game", middleware.Authentication(), middleware.GameState(), api.Game)
		//针对用户所选的数字调微信预支付
		apiRouter.POST("/prepay", middleware.Authentication(), middleware.GameState(), api.PrePay)
		//微信支付回调
		apiRouter.POST("/notify", api.Notify)
		// 开始游戏初始化
		apiRouter.GET("/play", middleware.Authentication(), middleware.GameState(), api.Play)
		// 用户砸蛋上报，返回是否砸中数字的结果
		apiRouter.POST("/smash", middleware.Authentication(), middleware.GameState(), api.Smash)
	}

	adminRouter := r.Group("/admin")
	{
		//提交后台登入
		adminRouter.POST("/login", admin.Login)
		//修改管理员密码
		adminRouter.POST("/password", middleware.AdminAuthentication(), admin.Password)
		//获取后台仪表盘数据
		adminRouter.GET("/", middleware.AdminAuthentication(), admin.Dashboard)
		//提交后台登出
		adminRouter.POST("/logout", middleware.AdminAuthentication(), admin.Logout)
		//获取登录验证码
		//adminRouter.GET("/captcha", middleware.AdminAuthentication(), admin.Captcha)
		//获取游戏设置
		adminRouter.GET("/show",  middleware.AdminAuthentication(), admin.Show)
		//修改游戏设置
		adminRouter.POST("/setting", middleware.AdminAuthentication(), admin.Setting)
		//重置游戏
		adminRouter.POST("/reset", middleware.AdminAuthentication(), admin.Reset)
		//开启游戏
		adminRouter.POST("/start", middleware.AdminAuthentication(), admin.Start)
		//停止游戏
		adminRouter.POST("/stop", middleware.AdminAuthentication(), admin.Stop)
	}

	r.GET("/test", func(c *gin.Context) {
		fmt.Println(time.Now())

	})

	r.GET("/init", func(c *gin.Context) {
		if err := service.Conn.Update(func(tx *nutsdb.Tx) error {
			err := tx.Delete("admin", []byte("admin"))
			if err != nil {
				log.Println(err)
			}
			password, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
			return tx.Put("admin", []byte("admin"), password,0)
		}); err != nil {
			log.Println(err)
		}
	})
}
