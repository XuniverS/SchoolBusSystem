package backend

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func RegisterSetupRoutes(router *gin.Engine) {
	busRoutes := router.Group("/bus")
	{
		busRoutes.POST("/addBus", addBus)
		busRoutes.POST("/deleteBus", removeBus)
		busRoutes.POST("/queryAll", queryAll)
		busRoutes.POST("/queryUser", queryUsersWithUserID)
	}
}

func addBus(c *gin.Context) {
	var requestBus struct {
		Origin      string `gorm:"column:origin" json:"origin"`
		Destination string `gorm:"column:destination" json:"destination"`
		BusType     string `gorm:"column:busType" json:"bustype"`
		Date        string `gorm:"column:date" json:"date"`
		Time        string `gorm:"column:time" json:"time"`
		Plate       string `gorm:"column:plate" json:"plate"`
		TotalSeats  int    `gorm:"column:total_seats" json:"totalseats"`
	}

	if err := c.ShouldBindJSON(&requestBus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	var bus Bus
	busDate, err := time.Parse("2006-01-02", requestBus.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "日期格式错误"})
		return
	}
	bus.Origin = requestBus.Origin
	bus.Destination = requestBus.Destination
	bus.BusType = requestBus.BusType
	bus.Date = busDate
	bus.Time = requestBus.Time
	bus.Plate = requestBus.Plate
	bus.TotalSeats = requestBus.TotalSeats

	if bus.BusType != "师生车" && bus.BusType != "教职工车" {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "班车类型错误！仅限为师生车或教职工车"})
		return
	}

	if err := insertBus(&bus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "数据库写入失败"})
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "添加成功"})
}

func removeBus(c *gin.Context) {
	var bus Bus
	if err := c.ShouldBindJSON(&bus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail"})
		return
	}
	result := deleteBus(bus.BusId)
	c.JSON(http.StatusOK, result)
}

func insertBus(bus *Bus) error {
	result := db.Create(bus)
	return result.Error
}

func deleteBus(busId int) map[string]interface{} {
	var bus Bus
	result := db.Where("busId =?", busId).Take(&bus)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return map[string]interface{}{"status": "fail", "message": "要删除的记录不存在"}
		}
		return map[string]interface{}{"status": "fail", "message": "查询记录时出错，无法执行删除操作"}
	}

	bookingCount := db.Model(&Bus{}).Where("busId =?", busId).Association("Bookings").Count()
	if bookingCount > 0 {
		return map[string]interface{}{"status": "booked", "message": "有用户预约了该班车，删除失败"}
	}

	deleteResult := db.Delete(&bus)
	if deleteResult.Error != nil {
		return map[string]interface{}{"status": "fail", "message": "删除记录时出错"}
	}
	if deleteResult.RowsAffected == 0 {
		return map[string]interface{}{"status": "fail", "message": "删除操作未影响任何记录，可能已被删除或其他原因"}
	}
	return map[string]interface{}{"status": "success", "message": "删除成功"}
}
