package api

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"smash-golden-eggs/service"
	"time"
)

type smashForm struct {
	Figure string `form:"figure" binding:"required"`
}

// Game 获取游戏信息
func Game(c *gin.Context) {
	payItems := service.GameInstance.Figures
	c.JSON(http.StatusOK, gin.H{
		"code":0,
		"msg": "选一个心仪的数字买了吧",
		"data": payItems,
	})
}

// Play 获取随机排序的payItems
func Play(c *gin.Context){
	payItems := service.GameInstance.Figures
	shuffle(payItems)

	c.JSON(http.StatusOK, gin.H{
		"code":0,
		"msg":"开始砸金蛋啦",
		"data":payItems,
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
	paidFigure := session.Get("paidFigure").(string)
	//todo 开启计时器倒计时，如果超过时间未重新支付，则解锁service.GameInstance.PlayMutex
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
	//todo 从Game.Figures中删除所砸的金蛋，将砸掉的金蛋序号写入Game.SmashedFigures中 （事务处理）

	//对比接口post上来的smash数字是否一致，如果一致返回成功砸中，不一致返回没砸中
	if paidFigure == c.PostForm("figure") {
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


func shuffle(slice []any) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(slice) > 0 {
		n := len(slice)
		randIndex := r.Intn(n)
		slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
		slice = slice[:n-1]
	}
}
