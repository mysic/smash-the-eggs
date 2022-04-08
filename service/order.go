package service

import (
	"github.com/bwmarrin/snowflake"
	"strconv"
)

const (
	OrderStateNoneExist = "nothing"
	OrderStatePaying = "paying"
	OrderStatePaid = "paid"
	OrderStateCancel = "cancel"
)

var OrderState string
func OrderSnGen() string {
	node, _ := snowflake.NewNode(1)
	return strconv.FormatInt(int64(node.Generate()), 10)
}


