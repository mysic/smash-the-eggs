package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/xujiajun/nutsdb"
	"net/http"
	"smash-golden-eggs/service"
	"strconv"
)

const (
	GameSettingBucket = "gameSetting"
	NumberSetting = "number"
	LockSecondLimit = "duration"
	Price = "price"
)

type settingForm struct {
	Number int  `form:"number" binding:"required"`
	LockSecondLimit int `form:"duration" binding:"required"`
	Price int `form:"price" binding:"required"`
}

// Show 显示游戏设置
func Show (c *gin.Context) {
	setting := map[string]string{}
	err := service.Conn.View(func(tx *nutsdb.Tx) error {
		bucket := GameSettingBucket
		all, err := tx.GetAll(bucket)
		if err != nil {
			return err
		}
		for _, item := range all {
			setting[string(item.Key)] = string(item.Value)
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":"查询错误",
			"data":err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":"查询成功",
		"data":setting,
	})
}

// Setting 更新游戏设置
func Setting(c *gin.Context) {
	var params settingForm
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":"参数错误",
			"data":err.Error(),
		})
		return
	}
	number := c.DefaultPostForm("number", "6")
	lockSecondLimit := c.DefaultPostForm("duration", "10")
	price := c.DefaultPostForm("price", "100")
	//更新数据库游戏设置
	err := service.Conn.Update(func(tx *nutsdb.Tx) error {
		bucket := GameSettingBucket
		key := []byte(NumberSetting)
		val := []byte(number)
		err := tx.Put(bucket, key, val, 0)
		if err != nil {
			return err
		}
		err = tx.Put(bucket, []byte(LockSecondLimit), []byte(lockSecondLimit), 0)
		if err != nil {
			return err
		}

		err = tx.Put(bucket, []byte(Price), []byte(price), 0)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":"设置失败",
			"data":err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":"设置成功",
	})
}

//Start 开始游戏
func Start(c *gin.Context) {
	err := gameInit(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":-1,
			"msg": "游戏开启失败，请联系管理员",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":0,
		"msg": "游戏已开启",
	})
}

func gameInit(c *gin.Context) error {
	figures,err := figureInit()
	if err != nil {
		return err
	}
	service.GameInstance.Figures = figures
	service.GameInstance.State = true
	return nil
}

func figureInit() ([]int64,error) {
	var figuresNumber uint64
	if err := service.Conn.View(func(tx *nutsdb.Tx) error {
		e,err := tx.Get(GameSettingBucket, []byte(NumberSetting))
		if err != nil {
			return err
		}

		figuresNumber, err = strconv.ParseUint(string(e.Value), 10, 8)
		if err != nil {
			return err
		}
		return nil
	});err != nil {
		return nil, err
	}
	var i uint64
	var figures []int64
	for i = 1;i <= figuresNumber; i++ {
		figures = append(figures, int64(i))
	}
	return figures, nil
}

//Stop 结束游戏
func Stop(c *gin.Context) {
	service.GameInstance.Figures =  []int64{}
	service.GameInstance.SmashedFigures = []int64{}
	service.GameInstance.CurrentPlayer = ""
	service.GameInstance.PayCount = 0
	service.GameInstance.State = false
	service.GameInstance.PlayMutex = false
	c.JSON(http.StatusOK, gin.H{
		"code":0,
		"msg": "游戏已停止",
	})
}
//Reset 重置游戏
func Reset(c *gin.Context) {
	service.GameInstance.Figures =  []int64{}
	service.GameInstance.SmashedFigures = []int64{}
	service.GameInstance.CurrentPlayer = ""
	service.GameInstance.PayCount = 0
	service.GameInstance.PlayMutex = false
	err := gameInit(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":-1,
			"msg": "游戏重置失败，请联系管理员",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":0,
		"msg": "游戏已重置",
	})
}
