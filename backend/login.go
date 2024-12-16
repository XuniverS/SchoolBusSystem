package backend

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterUserModule(router *gin.Engine) {

	userRouters := router.Group("/user")
	{
		userRouters.POST("/login", userLogin)
		userRouters.POST("/signin", userSignIn)
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
<<<<<<< HEAD
		fmt.Printf("shaEncode\n")
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail"})
=======
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "账号或密码错误"})
>>>>>>> origin/master
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "userid": queriedUser.UserID, "usertype": queriedUser.UserType, "isfirstlogin": updateUserIsFirstLogin(queriedUser)})
}

<<<<<<< HEAD
=======
func updateUserIsFirstLogin(user *User) int {
	if user.Is_first_login == 1 {
		db.Model(&User{}).Where("userId =?", user.UserID).Update("is_first_login", 0)
		return 1
	}
	return 0
}

>>>>>>> origin/master
func queryUser(user *User) (*User, error) {
	var queriedUser User
	result := db.Where("username = ?", user.UserName).Take(&queriedUser)
	if result.Error != nil {
		return &User{}, result.Error
	}
	return &queriedUser, nil
}

func updateUserIsFirstLogin(user *User) int {
	if user.Is_first_login {
		db.Model(&User{}).Where("userId =?", user.UserID).Update("is_first_login", 0)
		return 1
	}
	return 0
}

func userSignIn(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail"})
		return
	}

	if len(user.UserID) == 0 || len(user.UserName) == 0 || len(user.Password) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "userID、用户名、密码都不能为空"})
		return
	}
	user.Is_first_login = 1
	user.Password = shaEncode(user.Password)
	if err := insertUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "插入用户失败"})
		return
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

/*
func shaEncode(p string) string {
	hashBytes := sha256.Sum256([]byte(p))
	return hex.EncodeToString(hashBytes[:])
}
*/
