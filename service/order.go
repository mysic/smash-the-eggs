package service

import (
	"github.com/bwmarrin/snowflake"
	"strconv"
)

const (
	OrderStatusNoneExist = "nothing"
	OrderStatusPaying = "paying"
	OrderStatusPaid = "paid"
	OrderStatusCancel = "cancel"
)

var OrderStatus string

func OrderSnGen() string {
	node, _ := snowflake.NewNode(1)
	return strconv.FormatInt(int64(node.Generate()), 10)
}


