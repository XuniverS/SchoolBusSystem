package backend

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterSetupRoutes(router *gin.Engine) {
	busRoutes := router.Group("/bus")
	{
		busRoutes.POST("/addBus", addBus)
		busRoutes.POST("/deleteBus", removeBus)
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
