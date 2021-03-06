package api

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/xujiajun/nutsdb"
	"net/http"
	"smash-golden-eggs/admin"
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
	if 0 >= len(service.GameInstance.Figures) + len(service.GameInstance.SmashedFigures) {
		go func() {
			time.Sleep(time.Second * 1)
			admin.Reset(c)
		}()
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":"金蛋已经全部砸开了，请等待下一轮新游戏吧",
		})
		return
	}
	session := sessions.Default(c)
	//用户手机号
	mobile := session.Get("mobile")

	if (service.GameInstance.PlayMutex == true && mobile != service.GameInstance.CurrentPlayer) || !service.Mutex.TryLock() {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":"游戏正在进行中，请稍等片刻",
		})
		return
	}
	service.GameInstance.PlayMutex = true //tips 此变量的更新在最终游戏结束当前用户时限内未购买的情况下进行解锁
	service.Mutex.Unlock()
	// 用户打算购买的金蛋
	figure, _ := strconv.ParseInt(c.PostForm("figure"),0,0)
	if service.FindFigureInSlice(service.GameInstance.SmashedFigures, figure) >=0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":"这个金蛋已经砸过了，换一个砸吧",
		})
		return
	}

	service.OrderState = service.OrderStatePaying
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
	// 保存用户选的金蛋figure
	service.PaidFigure = figure
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
	//付款的金额
	payCount++
	payAmount := payCount * 100
	//生成订单,状态为支付中,保存入库
	//list - key订单号，val为手机号-支付金额-购买数字-状态-日期-时间戳
	//zSet - key订单号，范围查找在set中找，找到后从list中取详细信息
	bucket := service.OrderBucket
	key := []byte(orderSn)
	err = service.Conn.Update(func(tx *nutsdb.Tx) error {
		//手机号
		err := tx.LPush(bucket, key,
			[]byte(mobile.(string)),//手机号
			[]byte(strconv.FormatInt(int64(payAmount), 10)),//支付金额
			[]byte(strconv.FormatInt(figure,10)),//购买数字
			[]byte(service.OrderStatePaying),//购买状态
			[]byte(time.Now().String()),
			[]byte(strconv.FormatInt(time.Now().Unix(),10)),//时间戳
		)
		if err != nil {
			return err
		}
		err = tx.ZAdd(bucket, key, float64(time.Now().Unix()),[]byte(time.Now().String()))
		return nil
	})
	if err != nil {
		return 
	}
	// 拉起支付后，60秒内未完成支付，则关闭订单
	go func(bucket string, key []byte) {
		time.Sleep(time.Second * 60)
		err := service.Conn.View(func(tx *nutsdb.Tx) error {
			order, err := tx.LRange(bucket, key,0,-1)
			if err != nil {
				return err
			}
			if string(order[3]) != service.OrderStatePaid {
				service.GameInstance.PlayMutex = false
				//todo 调用微信关闭订单接口
			}

			return nil
		})
		if err != nil {
			return
		}
	}(bucket,key)
	//todo 获取商户支付参数,调用预下单接口(不返JSON，直接redirect)
	//client, err := service.NewWechatPayClient()
	//if err != nil {
	//	log.Fatalf("new wechat pay client err:%s", err)
	//}
	//goodsDetail := make([]h5.GoodsDetail,1)
	//goodsDetail = append(goodsDetail, h5.GoodsDetail{
	//	MerchantGoodsId:  core.String(figure),
	//	WechatpayGoodsId: nil,
	//	GoodsName:        core.String(figure),
	//	Quantity:         core.Int64(1),
	//	UnitPrice:        core.Int64(int64(payAmount * 100)),
	//})
	//ctx := context.Background()
	//wxApi := h5.H5ApiService{Client: client}
	//resp, result, err := wxApi.Prepay(ctx, h5.PrepayRequest{
	//	Appid:         core.String(service.AppID),
	//	Mchid:         core.String(service.MchID),
	//	Description:   core.String(mobile.(string) + "-" + figure + "-" + strconv.FormatInt(int64(payAmount), 10)),
	//	OutTradeNo:    core.String(orderSn),
	//	TimeExpire:    core.Time(time.Now()),
	//	Attach:        core.String("自定义数据说明"),
	//	NotifyUrl:     core.String(service.NotifyUrl),
	//	GoodsTag:      core.String(""),
	//	LimitPay:      make([]string, 1),
	//	SupportFapiao: core.Bool(false),
	//	Amount: &h5.Amount{
	//		Total: core.Int64(int64(payAmount)),
	//	},
	//	Detail: &h5.Detail{
	//		InvoiceId: core.String(orderSn),
	//		GoodsDetail: goodsDetail,
	//	},
	//	SceneInfo:&h5.SceneInfo{},
	//	SettleInfo: &h5.SettleInfo{},
	//})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//log.Println(resp)
	//log.Println(result)
	//fixme 临时测试，用完删除 {
	service.GameInstance.CurrentPlayer = mobile.(string)
	service.GameInstance.PayCount++
	//fixme }

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg" : "拉起支付成功",
		"data": orderSn,//todo 微信返回的预支付下单地址
	})
}

func Notify(c *gin.Context){
	var orderSn,resultCode string
	//fixme 测试用，用完删除 {
	orderSn =  c.PostForm("sn")
	resultCode = c.PostForm("result")
	//fixme }

	//todo 微信验签
	bucket := service.OrderBucket
	key := []byte(orderSn)
	if resultCode == "SUCCESS" {
		service.OrderState = service.OrderStatePaid
	} else {
		service.OrderState = service.OrderStateCancel
		service.GameInstance.PlayMutex = false
	}
	//更新订单状态
	err := service.Conn.Update(func(tx *nutsdb.Tx) error {
		if order, err := tx.LRange(bucket,key,0,-1); err != nil {
			return err
		} else {
			if string(order[3]) == service.OrderStatePaid {
				return errors.New(string(key) + "订单已支付")
			}
		}

		err := tx.LSet(bucket, key, 3, []byte(service.OrderState))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		c.String(http.StatusOK,err.Error())
		return
	}
	// 如果支付失败，解锁游戏，其他人可以购买
	if resultCode != "SUCCESS" {
		//ret := gin.H{
		//	"return_code" : "SUCCESS",
		//	"return_msg": "",
		//}
		// todo return response微信
		return
	}
	// 根据订单号查询所购买的figure
	var mobile string
	err = service.Conn.View(func(tx *nutsdb.Tx) error {
		order, err := tx.LRange(bucket, key, 0, -1)
		if err != nil {
			return err
		}
		mobile = string(order[5]) //val为手机号-支付金额-购买数字-状态-日期-时间戳
		return nil
	})
	if err != nil {
		return
	}
	// 订单记录下当前游戏中购买的
	service.GameInstance.CurrentPlayer = mobile
	service.GameInstance.PayCount++
	service.GameInstance.SmashPerm = true
	// todo return response微信
	//fixme 临时测试 {
	c.String(http.StatusOK,"success")
	//fixme }
}
