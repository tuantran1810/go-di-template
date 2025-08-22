# Go Dependency Injection Template Project

## Purpose
This is a **Go microservice template** that demonstrates Clean Architecture principles with dependency injection using Uber's fx framework. It serves as a starting point for building scalable, testable, and maintainable Go services.

## Key Features
- **Clean Architecture**: Structured in 4 layers (Entity, Use Case, Infrastructure, Interface Adapter)
- **Dependency Injection**: Uses Uber/fx for managing dependencies
- **gRPC + HTTP**: Dual protocol support with gRPC-Gateway for REST API generation
- **Generic Repository Pattern**: Reusable database patterns for multiple databases (MySQL, PostgreSQL, SQLite)
- **Data Transformation**: Generic transformers for converting between entities and DTOs
- **Comprehensive Testing**: Unit tests with mocks, integration tests with test containers
- **Protocol Buffers**: Uses Buf for proto file management and code generation
- **OpenTelemetry Integration**: Built-in observability with tracing, metrics, and structured logging

## Example Domain
The template implements a simple **User Management System** with:
- User creation and retrieval
- User attributes management
- Logging workers for asynchronous processing
- Message handling capabilities

## Tech Stack
- **Language**: Go 1.24.3
- **Framework**: Uber/fx for dependency injection
- **Database**: MySQL, PostgreSQL, SQLite support via GORM
- **API**: gRPC with gRPC-Gateway for HTTP REST
- **Testing**: testify + testcontainers for integration tests
- **Observability**: OpenTelemetry, Zap logger, Prometheus metrics
- **Tools**: Buf (protobuf), mockery (mocks), golangci-lint