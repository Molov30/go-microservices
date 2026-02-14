PROJECT_ROOT := $(abspath .)
GENERATED_DIR := $(PROJECT_ROOT)/generated
PROTO_DIR := $(PROJECT_ROOT)/proto
SERVICES_DIR := $(PROJECT_ROOT)/services

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')

.PHONY: gen help compose-up compose-down compose-logs compose-build compose-restart compose-ps compose-clean

## gen: Generate Go code from proto files
gen:
	$(MAKE) -C $(PROTO_DIR) gen GENERATED_DIR=$(GENERATED_DIR) PROJECT_ROOT=$(PROJECT_ROOT)

## compose-up: Start all services with docker compose
compose-up:
	@echo "Starting services with docker compose..."
	cd $(SERVICES_DIR) && VERSION=$(VERSION) BUILD_TIME=$(BUILD_TIME) docker compose up -d

## compose-down: Stop all services
compose-down:
	@echo "Stopping services..."
	cd $(SERVICES_DIR) && docker compose down

## compose-logs: View logs from all services
compose-logs:
	cd $(SERVICES_DIR) && docker compose logs -f

## compose-build: Build all services
compose-build:
	@echo "Building services..."
	cd $(SERVICES_DIR) && VERSION=$(VERSION) BUILD_TIME=$(BUILD_TIME) docker compose build

## compose-restart: Restart all services
compose-restart: compose-down compose-up

## compose-ps: List running services
compose-ps:
	cd $(SERVICES_DIR) && docker compose ps

## compose-clean: Stop services and remove volumes
compose-clean:
	@echo "Cleaning up services and volumes..."
	cd $(SERVICES_DIR) && docker compose down -v

## help: Print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
