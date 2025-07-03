package main

import (
	"Golang-CRUD/auth"
	"Golang-CRUD/internal/config"
	"Golang-CRUD/internal/infra/queue"
	httpDelivery "Golang-CRUD/delivery/http"
	"Golang-CRUD/internal/infra/repository"
	"Golang-CRUD/usecase"

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
	userRepo := repository.NewUserRepository(db)

	// Init usecase
	callUC := usecase.NewCallUsecase(callRepo)
	userUC := usecase.NewUserUsecase(userRepo)

	// Init auth handler (vẫn giữ)
	authHandler := auth.NewAuthHandler(db)

	// Init router từ file route.go
	router := httpDelivery.InitRouter(callUC, userUC, authHandler, db)

	// Run server
	router.Run(":8080")
}
