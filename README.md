# Update converse to Clean Architecture:
/domain        ← Core business logic (Entity, interface repository)
/usecase       ← Use cases (business logic)
/repository    ← Implement repository (DB, MQ), triển khai domain interface
/delivery      ← HTTP handler (Gin), gRPC, MQ consumer (adapter)
/infrastructure ← DB, MQ connection, 3rd party service
