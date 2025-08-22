# Development Tools and Workflow

## Core Development Tools

### Build and Dependency Management
- **Go Modules**: Dependency management (`go.mod`, `go.sum`)
- **Makefile**: Centralized build commands and workflows
- **Vendor**: Go module vendoring for reproducible builds

### Code Generation Tools
- **Buf**: Protocol buffer management and code generation
  - Configuration: `buf.yaml`, `buf.gen.yaml`, `buf.lock`
  - Generates Go gRPC code, gRPC-Gateway, and OpenAPI specs
- **Mockery**: Mock generation for testing
  - Configuration: `.mockery.yml`
  - Generates mocks in `mocks/` directory

### Code Quality Tools
- **golangci-lint**: Comprehensive Go linter
  - Configuration: `.golangci.yaml`
  - Runs extensive checks for code quality and style
- **gofmt/goimports**: Code formatting (automatic)

### Testing Tools
- **testify**: Testing framework with assertions and suites
- **testcontainers**: Integration testing with real databases
- **Go coverage**: Built-in coverage reporting

## Development Workflow

### 1. Initial Setup
```bash
make install-dev-env  # Install all required tools
make vendor          # Download dependencies
```

### 2. Code Generation (when needed)
```bash
make gen-proto       # After modifying .proto files
make gen-mock        # After modifying interfaces
```

### 3. Development Cycle
```bash
# Write code
make lint           # Check code quality
make test          # Run tests
make build         # Build application
```

### 4. Before Committing
```bash
make lint           # Ensure code quality
make test-coverage  # Verify test coverage
```

## Configuration Management

### Environment-based Configuration
- Uses `github.com/caarlos0/env/v11` for environment variable loading
- Structured configuration in `config/server_config.go`
- Support for nested configuration with prefixes

### Database Configuration
- Multiple database support (MySQL, PostgreSQL, SQLite)
- Connection pooling and migration support
- Test database via Docker Compose

## Protobuf Workflow

### Required Tools
- `protoc-gen-go`: Go protobuf compiler
- `protoc-gen-go-grpc`: gRPC Go compiler
- `protoc-gen-grpc-gateway`: HTTP/gRPC gateway
- `protoc-gen-openapiv2`: OpenAPI specification generation

### Generated Artifacts
- Go packages in `pkg/`
- OpenAPI specs in `swagger/`
- Automatic HTTP endpoint generation from gRPC services

## Testing Strategy

### Unit Tests
- Mock all external dependencies
- Table-driven test patterns
- Parallel test execution where possible
- High coverage requirements for business logic

### Integration Tests
- Real database instances via testcontainers
- End-to-end API testing
- Configuration and environment validation

## Debugging and Observability
- Structured logging with Zap
- Request correlation IDs
- OpenTelemetry integration for distributed tracing
- Prometheus metrics collection