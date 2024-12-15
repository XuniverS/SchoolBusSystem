package backend

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

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
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail"})
		return
	}
	queriedUser, err := queryUser(&user)
	if err != nil {
		fmt.Printf("queriedUser\n")
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail"})
		return
	}
	hashString := shaEncode(user.Password)
	if hashString != queriedUser.Password {
		fmt.Printf("shaEncode\n")
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "userid": queriedUser.UserID, "usertype": queriedUser.UserType, "isfirstlogin": updateUserIsFirstLogin(queriedUser)})
}

// 模拟的 updateUserIsFirstLogin 函数
func updateUserIsFirstLogin(user *User) int {
	if user.Is_first_login {
		return 1
	}
	return 0
}

// 模拟的 queryUser 函数
func queryUser(user *User) (*User, error) {
	// 这里返回一个固定的用户，模拟数据库查询
	if user.UserName == "user123" {
		return &User{
			UserID:         "1",
			UserType:       "admin",
			Password:       "hashed_password",
			Is_first_login: true,
			CreatedTime:    time.Now(),
		}, nil
	}
	return nil, errors.New("user not found")
}

/*func updateUserIsFirstLogin(user *User) int {
	if user.Is_first_login {
		db.Model(&User{}).Where("userId =?", user.UserID).Update("is_first_login", 0)
		return 1
	}
	return 0
}

func queryUser(user *User) (*User, error) {
	var queriedUser User
	result := db.Where("username = ?", user.UserName).Take(&queriedUser)
	if result.Error != nil {
		return &User{}, result.Error
	}
	return &queriedUser, nil
}*/

func userSignIn(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail"})
		return
	}

	if len(user.UserName) == 0 || len(user.Password) == 0 || len(user.Email) == 0 {
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

// 模拟的 shaEncode 函数
func shaEncode(password string) string {
	if password == "123456Aa" {
		return "hashed_password" // 模拟加密后的密码
	}
	return ""
}

/*func shaEncode(p string) string {
	hashBytes := sha256.Sum256([]byte(p))
	return hex.EncodeToString(hashBytes[:])
}*/
