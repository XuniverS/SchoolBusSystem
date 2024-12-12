package back

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	userID         string    `json:"userid"`
	userType       string    `json:"usertype"`
	userName       string    `json:"username"`
	email          string    `json:"email"`
	password       string    `json:"password"`
	if_first_login bool      `json:"if_first_login"`
	createdTime    time.Time `json:"createdtime"`
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
	user.userName = c.Query("username")
	user.password = c.Query("password")

	queriedUser, err := queryUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail"})
		return
	}
	hashBytes := sha256.Sum256([]byte(user.password))
	hashString := hex.EncodeToString(hashBytes[:])
	if hashString != queriedUser.password {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "userType": "学生", "init": "1"})
}

func queryUser(user User) (*User, error) {
	var queriedUser User
	result := db.Where("username = ?", user.userName).First(&queriedUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return &queriedUser, nil
}

func userSignIn(c *gin.Context) {
	var user User
	user.userType = c.Query("usertype")
	user.userName = c.Query("username")
	user.password = c.Query("password")
	user.email = c.Query("email")

	if user.userName == "" || user.password == "" || user.email == "" {
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
