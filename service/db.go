package service

import "github.com/xujiajun/nutsdb"

const (
	OrderBucket = "order"
	OrderSet = "orderSet"
	AdminBucket = "admin"
)

var Conn *nutsdb.DB