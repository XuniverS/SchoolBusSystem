package back

import (
	"crypto/sha256"
	"encoding/hex"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type User struct {
	userID         string    `json:"userId"`
	userType       string    `json:"userType"`
	userName       string    `json:"userName"`
	email          string    `json:"email"`
	password       string    `json:"password"`
	if_first_login bool      `json:"If_First_Login"`
	createdTime    time.Time `json:"createdTime"`
}

func LoginModule() *gin.Engine {
	router := gin.Default()

	router.Use(cors.Default())

	router.Static("", "./front")

	router.POST("/login", login)

	return router
}

func login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
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
	queriedUser, err := queryUser(db, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "查询错误"})
		return
	}
	hashBytes := sha256.Sum256([]byte(user.password))
	hashString := hex.EncodeToString(hashBytes[:])
	if hashString != queriedUser.password {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "密码错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "userType": "学生", "init": "1"})
}

func queryUser(db *gorm.DB, user User) (*User, error) {
	var queriedUser User
	result := db.Where("user_id =?", user.userName).First(&queriedUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return &queriedUser, nil
}

func SignInModule(c gin.Context) {

}
