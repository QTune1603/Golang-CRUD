package http

import (
	"Golang-CRUD/auth"
	"Golang-CRUD/usecase"
	"Golang-CRUD/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRouter(callUC *usecase.CallUsecase, userUC *usecase.UserUsecase, authHandler *auth.AuthHandler, db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	callHandler := NewCallHandler(callUC)
	userHandler := NewUserHandler(userUC)

	// Public routes
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.POST("/logout", middleware.JWTMiddleware(db), authHandler.Logout)

	// Protected routes
	authGroup := r.Group("/api")
	authGroup.Use(middleware.JWTMiddleware(db))
	{
		// Các route quản lý user
		authGroup.GET("/me", userHandler.GetMyInfo)
		authGroup.PUT("/me", userHandler.UpdateMyInfo)
		authGroup.DELETE("/me", userHandler.DeleteMyAccount)

		// Các route call API
		authGroup.POST("/calls", callHandler.Create)
		authGroup.GET("/calls", callHandler.List)
		authGroup.GET("/calls/:id", callHandler.GetByID)
		authGroup.PUT("/calls/:id", callHandler.Update)
		authGroup.DELETE("/calls/:id", callHandler.Delete)
	}

	return r
}
