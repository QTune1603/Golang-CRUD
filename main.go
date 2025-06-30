package main

import (
	"call-api/config"
	"call-api/consumer"
	"call-api/controller"
	"call-api/middleware"
	"call-api/model"
	"call-api/auth"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.InitDB()
	db.AutoMigrate(&model.CallLog{}, &model.User{}, &model.RevokedToken{})
	rabbitConn := config.InitRabbitMQ()
	defer rabbitConn.Close()

	go consumer.StartResultUpdater(rabbitConn, db)

	r := gin.Default()
	r.Use(middleware.CORSMiddleware()) 

	authHandler := auth.NewAuthHandler(db) 

	//Public routes
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.POST("/logout", middleware.JWTMiddleware(db), authHandler.Logout)

	// Protected routes
	authGroup := r.Group("/api")
	authGroup.Use(middleware.JWTMiddleware(db))
	{
		
		// authGroup.GET("/me", controller.GetMyInfo)
		authGroup.GET("/me", controller.GetMyInfo(db))
		authGroup.PUT("/me", controller.UpdateMyInfo(db))
		authGroup.DELETE("/me", controller.DeleteMyAccount(db))
	}

	//call API routes
	controller.SetupRoutes(r, db, rabbitConn)

	r.Run(":8080")
}

