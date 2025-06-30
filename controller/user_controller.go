package controller

import (
	"call-api/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Lấy thông tin người dùng hiện tại
func GetMyInfo(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("userID")

		var user model.User
		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"created_at": user.CreatedAt,
		})
	}
}

// Cập nhật username hoặc password
func UpdateMyInfo(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("userID")

		var input struct {
			Username *string `json:"username"`
			Password *string `json:"password"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		updateData := map[string]interface{}{}
		if input.Username != nil {
			updateData["username"] = *input.Username
		}
		if input.Password != nil {
			hash, _ := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
			updateData["password"] = string(hash)
		}

		if len(updateData) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nothing to update"})
			return
		}

		err := db.Model(&model.User{}).Where("id = ?", userID).Updates(updateData).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Updated"})
	}
}

// Xoá tài khoản người dùng
func DeleteMyAccount(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("userID")

		if err := db.Delete(&model.User{}, userID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Account deleted"})
	}
}
