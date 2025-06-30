package main

import (
	"call-api/auth"
	"call-api/config"
	"call-api/consumer"
	"call-api/delivery/http"
	"call-api/middleware"
	"call-api/repository"
	"call-api/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	// Init DB
	db := config.InitDB()
	db.AutoMigrate(&repository.CallModel{}, &repository.UserModel{}, &repository.RevokedToken{})


	// Init RabbitMQ
	rabbitConn := config.InitRabbitMQ()
	defer rabbitConn.Close()

	// Consumer chạy nền
	go consumer.StartResultUpdater(rabbitConn, db)

	// Init repository
	callRepo := repository.NewCallRepository(db)

	// Init usecase
	callUC := usecase.NewCallUsecase(callRepo)

	// Init handler
	callHandler := http.NewCallHandler(callUC)

	// Init router
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	// Auth handler (vẫn giữ)
	authHandler := auth.NewAuthHandler(db)

	// Public routes
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.POST("/logout", middleware.JWTMiddleware(db), authHandler.Logout)


	userRepo := repository.NewUserRepository(db)
	userUC := usecase.NewUserUsecase(userRepo)
	userHandler := http.NewUserHandler(userUC)

	// Protected routes
	authGroup := r.Group("/api")
	authGroup.Use(middleware.JWTMiddleware(db))
	{
		// Các route quản lý user
		authGroup.GET("/me", userHandler.GetMyInfo)
		authGroup.PUT("/me", userHandler.UpdateMyInfo)
		authGroup.DELETE("/me", userHandler.DeleteMyAccount)

		// Các route call API (có auth)
		authGroup.POST("/calls", callHandler.Create)
		authGroup.GET("/calls", callHandler.List)
		authGroup.GET("/calls/:id", callHandler.GetByID)
		authGroup.PUT("/calls/:id", callHandler.Update)
		authGroup.DELETE("/calls/:id", callHandler.Delete)
	}

	// Nếu muốn để call API public, thêm trực tiếp ngoài authGroup:
	// r.POST("/calls", callHandler.Create)
	// r.GET("/calls", callHandler.List)
	// ...

	// Run server
	r.Run(":8080")
}
