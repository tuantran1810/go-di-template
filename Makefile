BINARY_NAME=go-di-template

.PHONY: clean \
tidy \
vendor \
lint \
gen-proto \
build \
gen-mock \
test \
test-coverage \
test-coverage-html \
install-dev-env \
build-docker-dev-image \
run-docker-dev-env

clean:
	go clean

tidy:
	GOPROXY=https://proxy.golang.org go mod tidy

vendor: tidy
	GOPROXY=https://proxy.golang.org go mod vendor

lint: vendor
	GOFLAGS=-buildvcs=false golangci-lint run

gen-proto:
	buf dep update
	buf lint
	buf build
	buf generate

build: vendor
	go build -o ${BINARY_NAME} main.go

gen-mock: vendor gen-proto
	rm -rf mocks
	mockery

test: gen-mock
	go test ./internal/...

test-coverage: gen-mock
	go test ./internal/... ./libs/... -coverprofile=coverage.out

test-coverage-html: test-coverage
	go tool cover -html=coverage.out

install-dev-env:
	GOPROXY=https://proxy.golang.org go install github.com/vektra/mockery/v3@v3.2.5
	GOPROXY=https://proxy.golang.org go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.3.1
	GOPROXY=https://proxy.golang.org go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.1
	GOPROXY=https://proxy.golang.org go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	GOPROXY=https://proxy.golang.org go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.25.1
	GOPROXY=https://proxy.golang.org go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.25.1
	GOPROXY=https://proxy.golang.org GO111MODULE=on go install github.com/bufbuild/buf/cmd/buf@v1.48.0

build-docker-dev-image:
	docker build -t go-di-template:dev -f Dockerfile.dev .

run-docker-dev-env:
	docker run --rm -it -v $(PWD):/app -w /app go-di-template:dev
