package back

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Bus struct {
	BusId       int       `json:"busId"`
	Origin      string    `json:"origin"`
	Destination string    `json:"destination"`
	BusType     string    `json:"busType"`
	Date        time.Time `json:"date"`
	Time        time.Time `json:"time"`
	Plate       string    `json:"plate"`
	Seats       int       `json:"seats"`
}

func RegisterSetupRoutes(router *gin.Engine) {
	busRoutes := router.Group("/bus")
	{
		busRoutes.POST("/addBus", addBus)
	}
}

func addBus(c *gin.Context) {
	var bus Bus
	if err := c.ShouldBindJSON(&bus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	if bus.BusType != "师生车" && bus.BusType != "教职工车" {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "班车类型错误！仅限为师生车或教职工车"})
		return
	}
	if bus.Date.Before(time.Now()) {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "错误的时间！"})
		return
	}

	if err := insertBus(&bus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "数据库写入失败"})
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "添加成功"})
}

func removeBus(c *gin.Context) {

}

func insertBus(bus *Bus) error {
	result := db.Create(bus)
	return result.Error
}

func deleteBus(bus *Bus) error {
	return nil
}
