# Overview
- This document provides guidances for your development process

# Developemnt tools
### Mockery
- We use mockery for generating mock implementations
- All generated packages are stored in `mocks/` folder
- The mockery configuration is stored in `.mockery.yml`
- The basic configuration looks like:
```yaml
packages:
    this/is/the/firstpackage: # The package path
        config:
            dir: "mocks/firstpackage" # The mock directory
            filename: "mock_{{.InterfaceName}}.go" # Filename format
        interfaces:
            IDatabase: # Your Interface name here
                config:
            IHttpClient:
                config: # Your Interface name here
    this/is/the/secondpackage: # The package path
        config:
            dir: "mocks/secondpackage"
            filename: "mock_{{.InterfaceName}}.go"
        interfaces:
            IUsecase:
                config:
```

### Golang Linter
- We use golangci-lint to ensure good coding style and eliminate code smell.
- The linter configuration is stored in `.golangci.yaml`

### Buf build for protobuf generation
- We adopt buf as the main tool for generating gRPC code from proto files
- Proto files are stored in `proto/` folder
- Generated golang code packages are stored in `pkg/` folder
- Generated swagger json files are stored in `swagger/` folder
- Buf configuration files: `buf.gen.yaml`, `buf.lock`, `buf.yaml`
- Following protoc plugins must be installed:
    - protoc-gen-go
    - protoc-gen-go-grpc
    - protoc-gen-grpc-gateway (v2)
    - protoc-gen-openapiv2 (v2)

# Development operations
- The `Makefile` contains all needed commands for the developemnt process
- gvm (https://github.com/moovweb/gvm) is recommended for switching between go environments
- `Makefile` usage:
    - Install development environemnt:  `make install-dev-env` 
    - Clean up the build: `make clean`
    - Acquire go vendor: `make vendor`
    - Run linter: `make lint`
    - Generate golang code from proto files: `make gen-proto`
    - Build code: `make build`
    - Generate mock code: `make gen-mock`
    - Run unit tests: `make test`
    - Run unit tests with coverage reports: `make test-coverage`
    - Run unit tests with coverage reports by html: `test-coverage-html`
