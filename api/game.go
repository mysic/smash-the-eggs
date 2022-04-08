package api

import (
	"github.com/gin-gonic/gin"
	"log"
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
	if !service.GameInstance.SmashPerm {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":"目前不能开始游戏，是不是还没有支付呢？",
			"data":"",
		})
		return
	}
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
	//判断当前是否有砸金蛋的权限（1.支付获得权限）
	if !service.GameInstance.SmashPerm {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":"目前不能砸金蛋，是不是还没有支付呢？",
			"data":"",
		})
		return
	}

	var params smashForm
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":"参数错误",
			"data":err.Error(),
		})
		return
	}

	smashFigure,_ := strconv.ParseInt(c.PostForm("figure"),0,0)
	log.Println("paidFigure: " + strconv.FormatInt(service.PaidFigure,10))
	log.Println("smashFigure: " + strconv.FormatInt(smashFigure,10))
	//判断提交的数字是否是已经砸过的数字
	if service.FindFigureInSlice(service.GameInstance.SmashedFigures, smashFigure) >= 0 {
		c.JSON(http.StatusOK, gin.H{
			"code" : 0,
			"msg": strconv.FormatInt(smashFigure,10) + "已经砸过了",
			"data":"",
		})
		return
	}
	service.GameInstance.SmashPerm = false
	//从Game.Figures中删除所砸的金蛋
	service.GameInstance.Figures = service.RemoveSliceElement(service.GameInstance.Figures, smashFigure)
	//将砸掉的金蛋序号写入Game.SmashedFigures中
	service.GameInstance.SmashedFigures = append(service.GameInstance.SmashedFigures, smashFigure)
	//倒计时10秒，如果没有购买则解锁游戏
	//todo 当支付成功后的下一局游戏未到本次倒计时结束，也已结束，该如何处理
	// 当本次支付成功，那么应该传进来一个值，表示已经支付成功，然后子coroutine倒计时完成后不做任何处理
	go func(chan bool) {
		counter := make(chan bool,1)
		go func(chan bool) {
			time.Sleep(time.Second * 10)
			counter <- true
		}(counter)
		if <-service.PayState {

		}


		if <-counter {
			if service.OrderState != service.OrderStatePaid {
				service.GameInstance.PlayMutex = false
			}
		}
		defer close(counter)
	}(service.PayState)
	//对比接口post上来的smash数字是否一致，如果一致返回成功砸中，不一致返回没砸中
	if service.PaidFigure == smashFigure {
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
