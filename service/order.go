package service

import (
	"github.com/bwmarrin/snowflake"
	"strconv"
)

const (
	OrderStateNoneExist = "nothing" //订单不存在
	OrderStatePaying = "paying" //支付中
	OrderStatePaid = "paid" // 已支付
	OrderStateCancel = "cancel" //已取消
)

var OrderState string
func OrderSnGen() string {
	node, _ := snowflake.NewNode(1)
	return strconv.FormatInt(int64(node.Generate()), 10)
}


