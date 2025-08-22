# Task Completion Checklist

## After Making Code Changes

### 1. Code Generation (if applicable)
- [ ] If you modified `.proto` files → Run `make gen-proto`
- [ ] If you modified interfaces with mock configs → Run `make gen-mock`
- [ ] If you added new interfaces → Update `.mockery.yml` and run `make gen-mock`

### 2. Code Quality Checks
- [ ] Run `make lint` to check for linting issues
- [ ] Fix any linting errors or warnings
- [ ] Ensure code follows Go conventions and project style guide

### 3. Testing Requirements
- [ ] Write unit tests for new code (mandatory for use case layer)
- [ ] Run `make test` to execute all unit tests
- [ ] Ensure all tests pass
- [ ] Run `make test-coverage` to check coverage
- [ ] Aim for high test coverage, especially in business logic

### 4. Build Verification
- [ ] Run `make build` to ensure code compiles successfully
- [ ] Test the built binary if making significant changes

### 5. Dependencies
- [ ] If you added new Go dependencies → Run `make vendor`
- [ ] Check if new tools need to be added to `make install-dev-env`

## Before Committing to Git

### 1. Final Verification
- [ ] All tests pass: `make test`
- [ ] Linting passes: `make lint`
- [ ] Code builds successfully: `make build`
- [ ] Coverage is adequate: `make test-coverage`

### 2. Documentation
- [ ] Update comments and documentation if needed
- [ ] Ensure public functions have proper GoDoc comments
- [ ] Update relevant README or documentation files if needed

### 3. Git Workflow
- [ ] Review your changes: `git diff`
- [ ] Stage your changes: `git add .`
- [ ] Commit with descriptive message: `git commit -m "descriptive message"`
- [ ] Push to repository: `git push origin branch-name`

## Development Environment Checks

### Missing Tools Recovery
- [ ] If any tools are missing → Run `make install-dev-env`
- [ ] Verify all required protoc plugins are installed
- [ ] Check Go version compatibility (requires Go 1.24.3+)

### Database Setup (for integration tests)
- [ ] Start test database: `cd test && docker-compose up -d`
- [ ] Verify database connectivity
- [ ] Run integration tests if applicable
- [ ] Stop test database when done: `cd test && docker-compose down`

## Performance and Observability
- [ ] Check that new code includes proper error handling
- [ ] Ensure context propagation for timeouts and cancellation
- [ ] Add appropriate logging for debugging
- [ ] Consider OpenTelemetry tracing for new service boundaries
- [ ] Monitor resource usage for performance-critical code