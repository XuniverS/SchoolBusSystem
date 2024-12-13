package backend

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

type Booking struct {
	BookingId int       `gorm:"column:bookingId;primaryKey" json:"bookingid"`
	UserId    string    `gorm:"column:userId" json:"userid"`
	BusId     int       `gorm:"column:busId" json:"busid"`
	Status    string    `gorm:"column:status" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdtime"`
}

type User struct {
	UserID         string    `gorm:"column:userId;primaryKey" json:"userid"`
	UserType       string    `gorm:"column:userType" json:"usertype"`
	UserName       string    `gorm:"column:username" json:"username"`
	Email          string    `gorm:"column:Email" json:"email"`
	Password       string    `gorm:"column:password" json:"password"`
	Is_first_login bool      `gorm:"column:is_First_Login" json:"isfirstlogin"`
	CreatedTime    time.Time `gorm:"column:created_At" json:"createdtime"`
}

type Bus struct {
	BusId          int       `gorm:"column:busId;primaryKey" json:"busid"`
	Origin         string    `gorm:"column:origin" json:"origin"`
	Destination    string    `gorm:"column:destination" json:"destination"`
	BusType        string    `gorm:"column:busType" json:"bustype"`
	Date           time.Time `gorm:"column:date" json:"date"`
	Time           string    `gorm:"column:time" json:"time"`
	Plate          string    `gorm:"column:plate" json:"plate"`
	TotalSeats     int       `gorm:"column:total_seats" json:"totalseats"`
	AvailableSeats int       `gorm:"column:available_seats" json:"availableseats"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"createdtime"`
}

var db *gorm.DB

func InitDatabase() {
	var err error

	db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:C0137yx.@tcp(127.0.0.1:3306)/BusBookingSystem?parseTime=true",
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	sqlDB, err := db.DB()
	sqlDB.SetConnMaxIdleTime(time.Minute * 5)
	sqlDB.SetConnMaxLifetime(time.Hour * 2)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(20)
}
