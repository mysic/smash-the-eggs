package api

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/xujiajun/nutsdb"
	"net/http"
	"smash-golden-eggs/service"
	"strconv"
	"time"
)

type prePayForm struct {
	Figure int `form:"figure" binding:"required"`
}
//PrePay 支付预下单
func PrePay(c *gin.Context) {
	// todo 调用预下单的时候就锁定所选的figure
	// 游戏正在进行中，不能购买金蛋
	//todo 修改service.GameInstance.PlayMutex变量时，需要加互斥锁
	if service.GameInstance.PlayMutex == true {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":"游戏正在进行中，请稍等片刻,当前玩家["+service.GameInstance.CurrentPlayer+"]",
		})
		return
	}
	var params prePayForm
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":"参数错误",
			"data":err.Error(),
		})
		return
	}
	//判断金蛋是否已经全部购买完
	payCount := service.GameInstance.PayCount
	if payCount >= len(service.GameInstance.Figures) + len(service.GameInstance.SmashedFigures) {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":"金蛋已经全部砸开了，请等待下一轮新游戏吧",
		})
		return
	}
	// 用户买的金蛋保存到session
	figure := c.PostForm("figure")
	session := sessions.Default(c)
	session.Set("figure", figure)
	err := session.Save()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":"购买失败",
			"data":err.Error(),
		})
		return
	}
	//订单号
	orderSn := service.OrderSnGen()
	//用户手机号
	user := session.Get("mobile")
	//付款的金额
	payCount++
	payAmount := payCount * 100
	//生成订单,状态为支付中,保存入库
	//list - key订单号，val为手机号-支付金额-购买数字-状态-日期
	//zSet - key订单号，范围查找在set中找，找到后从list中取详细信息
	bucket := "order"
	key := []byte(orderSn)
	err = service.Conn.Update(func(tx *nutsdb.Tx) error {
		//手机号
		err := tx.LPush(bucket, key,
			[]byte(user.(string)),//手机号
			[]byte(strconv.FormatInt(int64(payAmount), 10)),//支付金额
			[]byte(figure),//购买数字
			[]byte(service.OrderStatusPaying),//购买状态
			[]byte(time.Now().String()),
			[]byte(strconv.FormatInt(time.Now().Unix(),10)),//时间戳
		)
		if err != nil {
			return err
		}
		err = tx.ZAdd(bucket, key, float64(time.Now().Unix()), key)
		return nil
	})
	if err != nil {
		return 
	}
	// todo 1分钟内未完成支付，则状态修改为可购买，可供其他人购买 修改状态通过goroutine计时器 + channel 完成
	//todo 获取商户支付参数,调用预下单接口
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg" : "用户购买成功",
	})
}

func Notify(c *gin.Context){
	//todo 微信验签
	//todo 根据订单号查询所购买的figure
	paidFigure := 1 //fixme temp delete
	session := sessions.Default(c)
	session.Set("paidFigure", paidFigure)
	//todo service:Game中存储本次支付成功的用户mobile
	//todo service.GameInstance.PayCount 加1
	//todo 订单记录中记录本次交易记录
	//todo 支付成功，开启互斥锁
}
