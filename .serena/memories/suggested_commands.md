# Essential Development Commands

## Environment Setup
```bash
# Install all development tools (run once)
make install-dev-env

# Clean and setup dependencies
make clean
make vendor
```

## Development Workflow
```bash
# Generate code from protobuf files
make gen-proto

# Generate mocks for testing
make gen-mock

# Run linter
make lint

# Build the application
make build
```

## Testing
```bash
# Run unit tests
make test

# Run tests with coverage
make test-coverage

# Generate HTML coverage report
make test-coverage-html
```

## Running the Application
```bash
# Run the server (gRPC + HTTP)
./go-di-template start-server

# Run the cron worker
./go-di-template start-cron

# Get help
./go-di-template --help
```

## Development Database
```bash
# Start test MySQL database
cd test && docker-compose up -d

# Stop test database
cd test && docker-compose down
```

## Useful macOS Commands
```bash
# Find Go files
find . -name "*.go" -type f

# Check file line counts
wc -l internal/*/*.go

# Search for text in files
grep -r "pattern" internal/

# List directory contents
ls -la

# Navigate directories
cd path/to/directory

# Git operations
git status
git add .
git commit -m "message"
git push origin main
```