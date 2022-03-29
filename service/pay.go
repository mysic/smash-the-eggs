package service

import (
	"context"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"log"
	"sync"
)

var (
	Mutex sync.Mutex
	MchID string = "1623464746"
	MchCertNumber string = ""
	MchAPIv3Key string = ""
	ApiClientKeyPath string = ""
	AppID string = ""
	NotifyUrl string = ""

)

func NewWechatPayClient() (*core.Client,error){
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(ApiClientKeyPath)
	if err != nil {
		log.Fatal("load merchant private key error")
	}

	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(MchID, MchCertNumber, mchPrivateKey, MchAPIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	return client,err
}