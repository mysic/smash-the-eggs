package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/xujiajun/nutsdb"
	"jindan/service"
	"net/http"
	"strconv"
)

// Dashboard 首页仪表盘
func Dashboard (c *gin.Context) {
	info := make(map[string]any)
	if len(service.GameInstance.Figures) > 0 {
		info["figures"]  = service.GameInstance.Figures
	} else {
		info["figures"] = []any{}
	}
	if service.GameInstance.SmashedFigures != nil {
		info["smashed"] = service.GameInstance.SmashedFigures
	} else {
		info["smashed"] = []any{}
	}
	info["status"] = service.GameInstance.Status
	info["player"] = service.GameInstance.CurrentPlayer
	c.JSON(http.StatusOK,gin.H{
		"code":0,
		"msg":"游戏信息",
		"data": info,
	})
}
//Start 开始游戏
func Start(c *gin.Context) {
	var figuresNumber uint64
	if err := service.Conn.View(func(tx *nutsdb.Tx) error {
		e,err := tx.Get(GameSettingBucket, []byte(NumberSetting))
		if err != nil {
			return err
		}

		figuresNumber, err = strconv.ParseUint(string(e.Value), 10, 8)
		if err != nil {
			return err
		}
		return nil
	});err != nil {
		return
	}
	var i uint64
	var figures []any
	for i = 1;i <= figuresNumber; i++ {
		figures = append(figures, i)
	}
	service.GameInstance.Figures = figures
	service.GameInstance.Status = true
	c.JSON(http.StatusOK, gin.H{
		"code":0,
		"msg": "游戏已开启",
	})
}
//Stop 结束游戏
func Stop(c *gin.Context) {
	service.GameInstance.Figures =  []any{}
	service.GameInstance.SmashedFigures = []any{}
	service.GameInstance.CurrentPlayer = ""
	service.GameInstance.Status = false
	c.JSON(http.StatusOK, gin.H{
		"code":0,
		"msg": "游戏已停止",
	})
}
//Reset 重置游戏
func Reset(c *gin.Context) {
	service.GameInstance.SmashedFigures = []interface{}{}
	service.GameInstance.CurrentPlayer = ""
	c.JSON(http.StatusOK, gin.H{
		"code":0,
		"msg": "游戏已重置",
	})
}