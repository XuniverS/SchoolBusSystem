package backend

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
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail"})
		return
	}

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
		"userType": queriedUser.UserType,
		"username": queriedUser.UserName,
		"email":    queriedUser.Email,
	})
}

func queryUserWithUserID(db *gorm.DB, user *User) (*User, error) {
	var queriedUser User
	result := db.Where("userId = ?", user.UserID).First(&queriedUser)
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
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail"})
		return
	}

	if user.UserID == "" || user.UserName == "" || user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Missing required fields"})
		return
	}

	existingUser, err := queryUserWithUserID(db, &user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		}
		return
	}

	existingUser.UserName = user.UserName
	existingUser.Email = user.Email

	if err := db.Save(&existingUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User information updated"})
}

func changePassword(c *gin.Context) {

}