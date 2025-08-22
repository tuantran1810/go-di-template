# Coding Style and Conventions

## Go Code Style
- **Formatting**: Enforced by `gofmt` and `goimports`
- **Linting**: Comprehensive linting with `golangci-lint`
- **Naming**: Go standard naming conventions (camelCase for private, PascalCase for public)
- **Error Handling**: Explicit error handling with wrapped errors using `fmt.Errorf("context: %w", err)`

## Architecture Rules

### Entity Layer Rules
- Must not depend on any other project packages
- Only import generic utility packages (time, strings, errors)
- Can use pure computation libraries (copier, mapstructure)
- Contains business objects as structs with simple operations

### Use Case Layer Rules
- Independent from infrastructure and frameworks
- All external dependencies abstracted behind interfaces
- Functions should be ≤ 100 lines (recommended)
- Must accept and return Entity objects
- Use context for timeouts and cancellation

### Infrastructure Layer Rules
- May depend on Entities and frameworks
- Implements interfaces defined in Use Case layer
- Contains data transformation logic
- Converts between Entity objects and external formats

## Testing Conventions
- **Unit Tests**: Table-driven test patterns with parallel execution
- **Mocking**: Generated mocks using mockery
- **Coverage**: Mandatory unit tests for Use Case layer
- **Integration Tests**: Use testcontainers for database testing
- **Test Structure**: Separate fast unit tests from slower integration tests

## Interface Design
- Small, purpose-specific interfaces
- Interface-driven development
- All public functions interact with interfaces, not concrete types
- Dependency inversion for external dependencies

## Error Handling
- Always check and handle errors explicitly
- Use wrapped errors for traceability
- Context-aware error messages
- Proper error types in entities package

## Documentation
- GoDoc-style comments for public functions and packages
- Clear, descriptive function names
- Document business logic and important decisions

## Observability
- Structured logging with Zap
- OpenTelemetry tracing for all service boundaries
- Context propagation for request correlation
- Prometheus metrics for key performance indicators