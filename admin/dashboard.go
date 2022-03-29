package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"smash-golden-eggs/service"
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
