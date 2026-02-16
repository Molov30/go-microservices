APP_NAME ?= account
BINARY_NAME ?= bin/$(APP_NAME)
MAIN_PATH ?= ./cmd
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS ?= -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

GOFILES ?= $(shell find . -name "*.go" -type f -not -path "./vendor/*" -not -path "./.git/*")

# Docker settings
DOCKER_IMAGE ?= $(APP_NAME)-service
DOCKER_TAG ?= $(VERSION)
DOCKER_REGISTRY ?= localhost:5000

# Default target
DEFAULT_GOAL := build

## build: Build the application
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p bin
	go build $(LDFLAGS) -o $(BINARY_NAME) $(MAIN_PATH)

## run: Run the application
run: build
	@echo "Running $(APP_NAME)..."
	./$(BINARY_NAME)

## run-only: Run without building (if binary exists)
run-only:
	./$(BINARY_NAME)

## deps: Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

## deps-update: Update dependencies
deps-update:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

## fmt: Format code
fmt:
	@echo "Formatting code..."
	gofmt -s -w $(GOFILES)
	gofumpt -l -w -s $(GOFILES)
	goimports -l -w $(GOFILES)

## lint: Run linter
lint:
	@echo "Running linter..."
	golangci-lint run --timeout=10m --fix -v ./...

## test: Run tests
test:
	@echo "Running tests..."
	go test -v ./...

## test-coverage: Run tests with coverage
test-coverage:
	@echo "Runnong tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## bench: Run benchmarks
bench:
	@echo "Running benchmarks..."
	go test -bench ./...

## clear: Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html

## install: Install the application
install: build
	@echo "Installing $(APP_NAME)"
	go install $(LDFLAGS) $(MAIN_PATH)

## tools: Install tools
tools:
	@echo "Installing tools..."
	go install go.uber.org/mock/mockgen@latest
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest

## docs: Generate documentation
docs:
	@echo "Generationg documentation..."
	go doc -http=:6060

## generate: Generate code from OpenAPI spec
generate:
	@echo "Generating code from OpenAPI spec..."
	@mkdir -p internal/generated
	oapi-codegen -package generated -generate types api/openapi.yml > internal/rest/v1/generated/types.gen.go
	oapi-codegen -package generated -generate std-http-server api/openapi.yml > internal/rest/v1/generated/server.gen.go


## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## dev: Development workflow
dev: fmt lint test build

## ci: CI/CD workflow
ci: deps fmt lint test build

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image $(DOCKER_IMAGE):$(DOCKER_TAG)..."
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		-t $(DOCKER_IMAGE):$(DOCKER_TAG) \
		-t $(DOCKER_IMAGE):latest \
		.

## docker-push: Push Docker image to registry
docker-push: docker-build
	@echo "Pushing Docker image to registry..."
	docker tag $(DOCKER_IMAGE):$(DOCKER_TAG) $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG)
	docker tag $(DOCKER_IMAGE):latest $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):latest
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG)
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):latest

## docker-run: Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run --rm -it \
		-p 8080:8080 \
		-p 50051:50051 \
		--name $(APP_NAME) \
		$(DOCKER_IMAGE):latest

## docker-clean: Remove Docker images
docker-clean:
	@echo "Removing Docker images..."
	docker rmi $(DOCKER_IMAGE):$(DOCKER_TAG) $(DOCKER_IMAGE):latest || true

.PHONY: build run run-nly deps deps-update fmt lint test test-coverage bench clean install tools docs help dev ci docker-build docker-push docker-run docker-clean
