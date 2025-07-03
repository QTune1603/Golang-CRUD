package main

import (
	"Golang-CRUD/auth"
	"Golang-CRUD/internal/config"
	"Golang-CRUD/internal/infra/repository"
	"Golang-CRUD/internal/infra/queue"
	httpDelivery "Golang-CRUD/delivery/http"
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
	go queue.StartResultUpdater(rabbitConn, db)

	// Init repositories
	callRepo := repository.NewCallRepository(db)
	userRepo := repository.NewUserRepository(db)
	readerRepo := repository.NewCallReaderRepository(db)

	// Init usecases
	callUC := usecase.NewCallUsecase(callRepo)
	userUC := usecase.NewUserUsecase(userRepo)

	// Init auth handler
	authHandler := auth.NewAuthHandler(db)

	// Init handlers
	callHandler := httpDelivery.NewCallHandler(callUC, readerRepo)
	userHandler := httpDelivery.NewUserHandler(userUC)

	// Init router từ file route.go
	router := httpDelivery.InitRouter(callHandler, userHandler, authHandler, db)

	// Run server
	router.Run(":8080")
}
