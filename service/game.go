package service

var GameInstance *Game
var AdminState bool
type Game struct {
	Figures []int64//生成的购买选项
	SmashedFigures []int64//已经砸开的选项
	CurrentPlayer string //当前游戏的玩家
	PayCount int //当前游戏已购买金蛋的次数
	Status bool // 游戏运行的状态 开始/结束
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



