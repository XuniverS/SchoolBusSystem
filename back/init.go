package back

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

func InitDatabase() {
	var err error
	// 这里只初始化一次数据库连接
	db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:C0137yx.@tcp(127.0.0.1:3306)/BusBookingSystem",
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
}
