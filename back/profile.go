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

	user.userID = c.DefaultQuery("userid", "")
	user.userName = c.DefaultQuery("username", "")
	user.email = c.DefaultQuery("email", "")

	if user.userID == "" || user.userName == "" || user.email == "" {
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

	existingUser.userName = user.userName
	existingUser.email = user.email

	if err := db.Save(&existingUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User information updated"})
}

func changePassword(c *gin.Context) {

}
