package back

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
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

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	router.Use(cors.Default())

	router.Static("", "./front")

	busRoutes := router.Group("/bus")
	{
		busRoutes.POST("/addBus", addBus)
	}

	return router
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
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:C0137yx.@tcp(127.0.0.1:3306)/BusBookingSystem",
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "数据库链接失败"})
		return
	}
	if err = insertBus(db, &bus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "数据库写入失败"})
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "添加成功"})
}

func insertBus(db *gorm.DB, bus *Bus) error {
	result := db.Create(bus)
	return result.Error
}
