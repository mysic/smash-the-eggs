package api

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"smash-golden-eggs/service"
	"strconv"
	"time"
)

type smashForm struct {
	Figure string `form:"figure" binding:"required"`
}

// Game 获取游戏信息
func Game(c *gin.Context) {
	data := make(map[string][]int64)
	data["figures"] = make([]int64,1)
	data["smashed_figures"] = make([]int64,1)
	data["figures"] = service.GameInstance.Figures
	data["smashed_figures"] = service.GameInstance.SmashedFigures
	c.JSON(http.StatusOK, gin.H{
		"code":0,
		"msg": "选一个心仪的数字买了吧",
		"data": data,
	})
}

// Play 获取随机排序的Figures
func Play(c *gin.Context){
	//todo 验证用户订单是否支付成功
	data := make(map[string][]int64)
	data["figures"] = make([]int64,1)
	data["smashed_figures"] = make([]int64,1)
	data["figures"] = service.GameInstance.Figures
	data["smashed_figures"] = service.GameInstance.SmashedFigures
	shuffle(data["figures"])
	c.JSON(http.StatusOK, gin.H{
		"code":0,
		"msg":"开始砸金蛋啦",
		"data":data,
	})
}

// Smash 验证用户砸中的payItem是否是他所购买的
func Smash(c *gin.Context) {
	var params smashForm
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":"参数错误",
			"data":err.Error(),
		})
		return
	}
	session := sessions.Default(c)
	paidFigure, _ := strconv.ParseInt(session.Get("figure").(string), 0,0)
	//todo 开启计时器倒计时，如果超过时间未重新支付，则解锁service.GameInstance.PlayMutex
	// 游戏结束10秒内未支付，则解锁。支付了则不解锁


	go func() {
		time.Sleep(time.Second * 10)
		//err := service.Conn.View(func(tx *nutsdb.Tx) error {
		//	orderInfo, err := tx.LRange(bucket, []byte(orderSn), 0, -1)
		//	if err != nil {
		//		return err
		//	}
		//	if string(orderInfo[3]) != service.OrderStatusPaid {
		//		service.GameInstance.PlayMutex = false
		//	}
		//
		//	return nil
		//})
		//if err != nil {
		//	return
		//}

	}()
	//从Game.Figures中删除所砸的金蛋
	smashFigure,_ := strconv.ParseInt(c.PostForm("figure"),0,0)
	service.GameInstance.Figures = service.RemoveSliceElement(service.GameInstance.Figures, smashFigure)
	//将砸掉的金蛋序号写入Game.SmashedFigures中
	service.GameInstance.SmashedFigures = append(service.GameInstance.SmashedFigures, smashFigure)

	//对比接口post上来的smash数字是否一致，如果一致返回成功砸中，不一致返回没砸中
	if paidFigure == smashFigure {
		c.JSON(http.StatusOK, gin.H{
			"code" : 0,
			"msg":"砸中啦",
			"data":"success",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code" : 0,
		"msg":"未砸中",
		"data":"fail",
	})

}


func shuffle(slice []int64) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(slice) > 0 {
		n := len(slice)
		randIndex := r.Intn(n)
		slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
		slice = slice[:n-1]
	}
}
