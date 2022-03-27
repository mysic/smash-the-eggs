package service

import (
	"github.com/bwmarrin/snowflake"
	"strconv"
)

const (
	OrderStatusPaying = "paying"
	OrderStatusPaid = "paid"
)

func OrderSnGen() string {
	node, _ := snowflake.NewNode(1)
	return strconv.FormatInt(int64(node.Generate()), 10)
}
