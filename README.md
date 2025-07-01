# Converse to Clean Architecture:
- /domain        ← Core business logic (Entity, interface repository)
  Chứa entity (model chuẩn, không phụ thuộc framework).
- /usecase       ← Use cases (business logic)
  Business logic thuần, chỉ phụ thuộc domain (không phụ thuộc hạ tầng).
- /repository    ← Implement repository (DB, MQ), triển khai domain interface
  Triển khai (implementation) data access, phụ thuộc database (MySQL).
- /delivery      ← HTTP handler (Gin), gRPC, MQ consumer (adapter)
  Các handler HTTP (Gin) để expose API.
- /middleware: Xử lý JWT, revoke, CORS.

- /consumer: Consumer xử lý message từ RabbitMQ.
