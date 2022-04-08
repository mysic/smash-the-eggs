package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/xujiajun/nutsdb"
	"log"
	"os"
	"smash-golden-eggs/middleware"
	"smash-golden-eggs/service"
)

const (
	LogFile  = "data/logs/api.log"
)


func main(){
	//fixme 发布时用release模式
	//gin.SetMode(gin.ReleaseMode)
	var err error
	logFile, _ := os.OpenFile(LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	router := gin.Default()
	// 启动session
	store := cookie.NewStore([]byte("jfa8389sf729ap(*&DJA(#xl"))
	router.Use(sessions.Sessions("smashGoldenEggs", store), middleware.Cors())
	service.AdminState = false
	// 启动db
	dbConfig := nutsdb.DefaultOptions
	dbConfig.Dir = "data/db"
	service.Conn, err = nutsdb.Open(dbConfig)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer func(Db *nutsdb.DB) {
		err := Db.Close()
		if err != nil {

		}
	}(service.Conn)
	// 注册路由
	RegisterRouter(router)
	// 实例化游戏
	service.GameInstance = &service.Game{
		Figures: []int64{},
		SmashedFigures: []int64{},
		CurrentPlayer: "",
		PayCount:0,
		State: false,
		PlayMutex: false,
	}
	//初始化订单状态
	service.OrderState = service.OrderStateNoneExist //订单未创建
	// 启动http服务
	err = router.Run(":8668")
	if err != nil {
		log.Println(err.Error())
	}
}