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
	//判断金蛋是否已经全部购买完
	payCount := service.GameInstance.PayCount
	if payCount >= len(service.GameInstance.Figures) + len(service.GameInstance.SmashedFigures) {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":"金蛋已经全部砸开了，请等待下一轮新游戏吧",
		})
		return
	}
	session := sessions.Default(c)
	mobile := session.Get("mobile")
	if (service.GameInstance.PlayMutex == true && mobile != service.GameInstance.CurrentPlayer) || !service.Mutex.TryLock() {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":"游戏正在进行中，请稍等片刻",
		})
		return
	}
	service.GameInstance.PlayMutex = true //tips 此变量的更新在最终游戏结束，当前用户时限内未购买的情况下进行解锁
	service.Mutex.Unlock()
	service.GameInstance.CurrentPlayer = mobile.(string)
	var params prePayForm
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":"参数错误",
			"data":err.Error(),
		})
		return
	}
	// 用户选的金蛋figure保存到session
	figure := c.PostForm("figure")
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
	//todo
	go func(bucket,orderSn string) {
		time.Sleep(time.Second * 60)
		err := service.Conn.View(func(tx *nutsdb.Tx) error {
			orderInfo, err := tx.LRange(bucket, []byte(orderSn), 0, -1)
			if err != nil {
				return err
			}
			if string(orderInfo[3]) != service.OrderStatusPaid {
				service.GameInstance.PlayMutex = false
			}

			return nil
		})
		if err != nil {
			return
		}

	}(bucket,orderSn)
	//todo 获取商户支付参数,调用预下单接口(不返JSON，直接redirect)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg" : "拉起支付成功",
		"data": "",//todo 微信返回的预支付下单地址
	})
}

func Notify(c *gin.Context){
	var orderSn,resultCode string
	orderSn = "32393932902039"//fixme 模拟微信回调中的 out_trade_no
	//resultCode := "success" //fixme 模拟微信回调中的结果
	//todo 微信验签
	bucket := "order"
	key := []byte(orderSn)
	payState := service.OrderStatusPaying
	if resultCode == "SUCCESS" {
		payState = service.OrderStatusPaid
	} else {
		payState = service.OrderStatusCancel
	}
	//更新订单状态
	err := service.Conn.Update(func(tx *nutsdb.Tx) error {
		err := tx.LSet(bucket, key, 3, []byte(payState))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return
	}
	// 如果支付失败，解锁游戏，其他认可以购买
	if resultCode != "SUCCESS" {
		service.GameInstance.PlayMutex = false
		//ret := gin.H{
		//	"return_code" : "SUCCESS",
		//	"return_msg": "",
		//}
		// todo return response微信
		return
	}
	// 根据订单号查询所购买的figure
	var paidFigure,mobile string
	err = service.Conn.View(func(tx *nutsdb.Tx) error {
		orderInfo, err := tx.LRange(bucket, key, 0, -1)
		if err != nil {
			return err
		}
		paidFigure = string(orderInfo[2])
		return nil
	})
	if err != nil {
		return
	}
	// 订单记录下当前游戏中购买的
	service.OrderInstance.PaidFigure = paidFigure
	service.GameInstance.CurrentPlayer = mobile
	service.GameInstance.PayCount++
	// todo return response微信

}
