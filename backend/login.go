package backend

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	UserID         string    `gorm:"column:userId;primaryKey" json:"userid"`
	UserType       string    `gorm:"column:userType" json:"usertype"`
	UserName       string    `gorm:"column:username" json:"username"`
	Email          string    `gorm:"column:email" json:"email"`
	Password       string    `gorm:"column:password" json:"password"`
	Is_first_login bool      `gorm:"column:is_first_login" json:"isfirstlogin"`
	CreatedTime    time.Time `gorm:"column:created_at" json:"createdtime"`
}

func RegisterUserModule(router *gin.Engine) {

	userRouters := router.Group("/user")
	{
		userRouters.POST("/login", userLogin)
		userRouters.GET("/signin", userSignIn)
	}

}

func userLogin(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail"})
		return
	}

	queriedUser, err := queryUser(user)
	fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail"})
		return
	}
	hashBytes := sha256.Sum256([]byte(user.Password))
	hashString := hex.EncodeToString(hashBytes[:])
	if hashString != queriedUser.Password {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "userType": "学生", "init": "1"})
}

func queryUser(user User) (*User, error) {
	var queriedUser User
	result := db.Where("username = ?", user.UserName).Take(&queriedUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return &queriedUser, nil
}

func userSignIn(c *gin.Context) {
	var user User
	user.UserType = c.Query("usertype")
	user.UserName = c.Query("username")
	user.Password = c.Query("password")
	user.Email = c.Query("email")

	if user.UserName == "" || user.Password == "" || user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "用户名、密码、邮箱都不能为空"})
		return
	}

	if insertUser(&user) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "插入用户失败"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户注册成功"})
}

func insertUser(user *User) error {
	result := db.Create(user)
	return result.Error
}
