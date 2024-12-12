package back

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func RegisterProfileModule(router *gin.Engine) {
	profilerouter := router.Group("/profile")
	{
		profilerouter.POST("/queryUser", queryUsersWithUserID)
		profilerouter.POST("/submitUser", submitUserInfo)
		profilerouter.POST("/changePassword", changePassword)
	}
}

func queryUsersWithUserID(c *gin.Context) {
	var user User
	user.userID = c.Query("userid")

	queriedUser, err := queryUserWithUserID(db, &user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"status": "not_found", "message": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		}
		return
	}

	// 成功返回用户信息
	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"userType": queriedUser.userType,
		"username": queriedUser.userName,
		"email":    queriedUser.email,
	})
}

func queryUserWithUserID(db *gorm.DB, user *User) (*User, error) {
	var queriedUser User
	result := db.Where("userId = ?", user.userID).First(&queriedUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}
	return &queriedUser, nil
}

func submitUserInfo(c *gin.Context) {
	var user User

	// 获取查询参数并赋值到结构体字段
	user.userID = c.DefaultQuery("userid", "")     // 如果没有提供userId，默认空字符串
	user.userName = c.DefaultQuery("username", "") // 如果没有提供username，默认空字符串
	user.email = c.DefaultQuery("email", "")       // 如果没有提供email，默认空字符串

	// 检查是否有必要的字段
	if user.userID == "" || user.userName == "" || user.email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Missing required fields"})
		return
	}

	// 查询数据库中是否存在该用户
	existingUser, err := queryUserWithUserID(db, &user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果找不到用户，返回404
			c.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "User not found"})
		} else {
			// 其他数据库错误
			c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		}
		return
	}

	// 更新用户信息
	existingUser.userName = user.userName
	existingUser.email = user.email

	// 保存更新的用户信息到数据库
	if err := db.Save(&existingUser).Error; err != nil {
		// 如果保存失败，返回修改失败
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Failed to update user"})
		return
	}

	// 如果更新成功，返回成功状态
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User information updated"})
}

func changePassword(c *gin.Context) {

}
