package service

var GameInstance *Game
var PaidFigure int64

type Game struct {
	Figures []int64//生成的购买选项
	SmashedFigures []int64//已经砸开的选项
	CurrentPlayer string //当前游戏的玩家
	PayCount int //当前游戏已购买金蛋的次数
	State bool // 游戏运行的状态 开始/结束
	SmashPerm bool // 是否允许砸金蛋 （未支付状态不能砸，已支付状态可以砸）
	PlayMutex bool //游戏是否加锁状态
}

func RemoveSliceElement(slice []int64, elem int64) []int64{
	if len(slice) == 0 {
		return slice
	}
	for i, v := range slice {
		if v == elem {
			slice = append(slice[:i], slice[i+1:]...)
			break
			}
		}
	return slice
}

func FindFigureInSlice(slice []int64, val int64) int {
	for i, item := range slice {
		if item == val {
			return i
		}
	}
	return -1
}



