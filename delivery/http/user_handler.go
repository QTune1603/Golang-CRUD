package http

import (
	"Golang-CRUD/domain"
	"Golang-CRUD/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	Usecase *usecase.UserUsecase
}

func NewUserHandler(uc *usecase.UserUsecase) *UserHandler {
	return &UserHandler{Usecase: uc}
}

func (h *UserHandler) GetMyInfo(c *gin.Context) {
	userID := c.GetUint("userID")
	user, err := h.Usecase.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":        user.ID,
		"username":  user.Username,
		"created_at": user.CreatedAt,
	})
}

func (h *UserHandler) UpdateMyInfo(c *gin.Context) {
	userID := c.GetUint("userID")
	var input domain.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Usecase.Update(userID, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated"})
}

func (h *UserHandler) DeleteMyAccount(c *gin.Context) {
	userID := c.GetUint("userID")
	if err := h.Usecase.Delete(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Account deleted"})
}
