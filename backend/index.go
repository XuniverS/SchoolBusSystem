package backend

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func RegisterIndexModule(router *gin.Engine) {
	indexRouter := router.Group("/index")
	{
		indexRouter.POST("/queryAll", queryAll)
		indexRouter.POST("/book", book)
		indexRouter.POST("/payed", payed)
		indexRouter.POST("/unbook", unbook)
		indexRouter.POST("/queryBooked", queryBooked)
		indexRouter.POST("/queryFinished", queryFinished)
	}
}

func queryAll(c *gin.Context) {
	var reqData struct {
		Date     string `json:"date"`
		UserType string `json:"usertype"`
	}
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "参数解析错误"})
		return
	}

	var buses []Bus
	if reqData.Date == "any" && reqData.UserType == "admin" {
		result := db.Find(&buses)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "any日期格式错误"})
			return
		}
		c.JSON(http.StatusOK, buses)
		return
	}
	queryDate, err := time.Parse("2006-01-02", reqData.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "日期格式错误"})
		return
	}
	var busType string
	if reqData.UserType == "学生" {
		busType = "师生车"
	} else {
		busType = "教职工车"
	}

	result := db.Where("date =? AND busType LIKE?", queryDate, "%"+busType+"%").Find(&buses)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, buses)
}

func book(c *gin.Context) {
	var reqData struct {
		UserId string `json:"userid"`
		BusId  int    `json:"busid"`
	}
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "参数解析错误"})
		return
	}

	var existingBooking Booking
	result := db.Where("userId =? AND busId =? AND status =?", reqData.UserId, reqData.BusId, "已预约").Take(&existingBooking)
	if result.Error == nil {
		c.JSON(http.StatusOK, gin.H{"status": "booked"})
		return
	}

	var bus Bus
	result = db.Where("busId =?", reqData.BusId).Take(&bus)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "查询班车信息失败"})
		return
	}

	if bus.AvailableSeats <= 0 {
		c.JSON(http.StatusOK, gin.H{"status": "full"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "available"})
}

func payed(c *gin.Context) {
	var reqData struct {
		UserId string `json:"userid"`
		BusId  int    `gorm:"column:busId" json:"busid"`
	}
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "参数解析错误"})
		return
	}
	var bus Bus
	result := db.Where("busId =?", reqData.BusId).Take(&bus)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "查询班车信息失败"})
		return
	}

	bus.AvailableSeats--
	if err := updateBus(&bus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "更新班车座位信息失败"})
		return
	}

	newBooking := Booking{
		UserId: reqData.UserId,
		BusId:  reqData.BusId,
		Status: "已预约",
	}
	if err := insertBooking(&newBooking); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "预约失败"})
		return
	}

	result = db.Save(&newBooking)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "更新预约记录状态失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
func unbook(c *gin.Context) {
	var reqData struct {
		UserId string `json:"userId"`
		BusId  int    `gorm:"column:busId" json:"busid"`
	}
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "参数解析错误"})
		return
	}

	var a Booking
	result := db.Where("userId =? AND busId =?", reqData.UserId, reqData.BusId).Take(&a)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "查询预约记录失败"})
		return
	}
	var booking Booking
	result = db.Where("userId =? AND busId =? AND status =?", a.UserId, a.BusId, "已预约").Take(&booking)

	if booking.Status != "已预约" && booking.Status != "已支付" {
		c.JSON(http.StatusOK, gin.H{"status": "fail", "message": "该预约记录不符合取消条件"})
		return
	}

	booking.Status = "已取消"
	result = db.Save(&booking)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "更新预约记录状态失败"})
		return
	}

	var bus Bus
	result = db.Where("busId =?", reqData.BusId).Take(&bus)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "查询班车信息失败"})
		return
	}
	bus.AvailableSeats++
	if err := updateBus(&bus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "更新班车座位信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func queryBooked(c *gin.Context) {
	var reqData struct {
		UserID string `json:"userid"`
	}
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "参数解析错误"})
		return
	}

	var bookings []Booking
	result := db.Where("userId =? AND (status = '已预约' OR status = '已支付')", reqData.UserID).Find(&bookings)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "查询预约记录失败"})
		return
	}

	var bookedBuses []Bus
	for _, booking := range bookings {
		var bus Bus
		result = db.Where("busId =?", booking.BusId).Take(&bus)
		if result.Error == nil {
			bookedBuses = append(bookedBuses, bus)
		}
	}
	if bookedBuses == nil {
		c.JSON(http.StatusOK, nil)
		return
	}

	c.JSON(http.StatusOK, bookedBuses)
}

func queryFinished(c *gin.Context) {
	var reqData struct {
		UserId string `json:"userId"`
	}
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "参数解析错误"})
		return
	}

	var bookings []Booking
	result := db.Where("userId =? AND (status = '已完成' OR status = '已取消')", reqData.UserId).Find(&bookings)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "查询预约记录失败"})
		return
	}

	var finishedBuses []struct {
		Status      string    `json:"status"`
		BusId       int       `json:"busid"`
		Origin      string    `json:"origin"`
		Destination string    `json:"destination"`
		Time        string    `json:"time"`
		BusType     string    `json:"bustype"`
		Plate       string    `json:"plate"`
		Date        time.Time `json:"date"`
	}
	for _, booking := range bookings {
		var bus Bus
		result = db.Where("busId =?", booking.BusId).Take(&bus)
		if result.Error == nil {
			var item struct {
				Status      string    `json:"status"`
				BusId       int       `json:"busid"`
				Origin      string    `json:"origin"`
				Destination string    `json:"destination"`
				Time        string    `json:"time"`
				BusType     string    `json:"bustype"`
				Plate       string    `json:"plate"`
				Date        time.Time `json:"date"`
			}
			if booking.Status == "已完成" {
				item.Status = "finished"
			} else {
				item.Status = "unbooked"
			}
			item.BusId = bus.BusId
			item.Origin = bus.Origin
			item.Destination = bus.Destination
			item.Time = bus.Time
			item.BusType = bus.BusType
			item.Plate = bus.Plate
			item.Date = bus.Date
			finishedBuses = append(finishedBuses, item)
		}
	}
	if finishedBuses == nil {
		c.JSON(http.StatusOK, Bus{})
	}
	c.JSON(http.StatusOK, finishedBuses)
}

func insertBooking(booking *Booking) error {
	result := db.Create(booking)
	return result.Error
}

func updateBus(bus *Bus) error {
	result := db.Save(bus)
	return result.Error
}
