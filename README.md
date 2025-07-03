# ✨ Converse to Clean Architecture

```plaintext
Golang-CRUD-main/
├── .env.example
├── .gitignore
├── Dockerfile
├── README.md
├── docker-compose.yml
├── go.mod
├── go.sum
├── auth/
│   └── auth_controller.go
├── cmd/
│   └── main.go
├── delivery/
│   └── http/
│       ├── call_handler.go
│       ├── route.go
│       └── user_handler.go
├── domain/
│   ├── call.go
│   └── user.go
├── internal/
│   ├── config/
│   │   ├── db.go
│   │   └── rabbitmq.go
│   ├── infra/
│   │   ├── queue/
│   │   │   └── result_updater.go
│   │   └── repository/
│   │       ├── call_log_repository.go
│   │       ├── call_reader_repository.go
│   │       ├── revoked_token_repository.go
│   │       └── user_repository.go
│   └── middleware/
│       ├── auth.go
│       └── cors.go
├── usecase/
│   ├── call_usecase.go
│   ├── user_usecase.go
│   └── reader/
│       └── call_reader_repository.go

Domain ← Usecase ← Delivery (Interface Adapter) ← Infra (Repository, Queue) ← External (DB, RabbitMQ)
