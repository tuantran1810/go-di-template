# Project Structure and Clean Architecture

## Directory Structure
```
go-di-template/
├── cmd/                    # Application entrypoints (CLI commands)
├── config/                 # Configuration schemas and loading
├── internal/               # Core application logic (private)
│   ├── entities/          # Domain entities (business objects)
│   ├── usecases/          # Business logic layer
│   ├── controllers/       # HTTP/gRPC controllers
│   ├── repositories/      # Data access layer
│   ├── inbound/           # Message consumers
│   └── outbound/          # External clients
├── libs/                  # Internal utilities
│   ├── logger/            # Logging utilities
│   ├── middlewares/       # HTTP/gRPC middlewares
│   ├── server/            # Server setup
│   └── utils/             # Common utilities
├── mocks/                 # Generated mock files
├── pkg/                   # Generated protobuf code
├── proto/                 # Protocol buffer definitions
├── swagger/               # Generated OpenAPI specs
├── test/                  # Test utilities and setup
└── vendor/                # Go module dependencies
```

## Clean Architecture Layers

### 1. Entity Layer (`internal/entities/`)
- Core business entities and domain objects
- No dependencies on other layers
- Contains data transformers for type conversion
- Examples: User, UserAttribute, Transaction

### 2. Use Case Layer (`internal/usecases/`)
- Application-specific business logic
- Orchestrates data flow between entities and infrastructure
- Depends only on entities and abstracted interfaces
- Examples: Users (user management), LoggingWorker

### 3. Infrastructure Layer
- **Controllers** (`internal/controllers/`): gRPC/HTTP API handlers
- **Repositories** (`internal/repositories/`): Database access patterns
- **Inbound** (`internal/inbound/`): Message consumers
- **Outbound** (`internal/outbound/`): External service clients

### 4. Interface Adapter Layer
- Embedded within Infrastructure layer
- Handles data transformations between layers
- Located in `internal/controllers/transformers/`

## Key Design Patterns

### Dependency Injection (Uber/fx)
- All dependencies injected via constructors
- Configured in DI containers
- Promotes testability and modularity

### Generic Repository Pattern
- Reusable database operations
- Supports MySQL, PostgreSQL, SQLite
- Type-safe generic implementations

### Data Transformation
- Base transformers using copier library
- Extended transformers for array operations
- Centralized conversion logic