PROTO_ROOT := .
DOCKER_IMAGE := proto-builder
GENERATED_DIR ?= $(abspath $(PROTO_ROOT)/../generated)
PROJECT_ROOT ?= $(abspath $(PROTO_ROOT)/..)

# Default target
DEFAULT_GOAL := gen

.PHONY: docker-build gen clean help

## docker-build: Build Docker image for proto generation
docker-build:
	docker build -t $(DOCKER_IMAGE) .

## gen: Generate Go code from proto files
gen: docker-build
	@echo "Generating Go code from proto files..."
	@mkdir -p $(GENERATED_DIR)
	@docker run --rm \
		--user $(shell id -u):$(shell id -g) \
		-v $(abspath $(PROTO_ROOT)):/app/proto \
		-v $(GENERATED_DIR):/app/generated \
		$(DOCKER_IMAGE) \
    	bash -c '\
    		set -e; \
    		for dir in account pagination auth gateway; do \
    			echo "Processing $$dir"; \
    			mkdir -p /app/generated/$$dir; \
    			cd /app/proto/$$dir; \
    			for file in *.proto; do \
    				echo "    Generating $$dir/$$file"; \
    				protoc \
    				    -I /app/proto/$$dir \
    					-I /app/proto \
    					-I /usr/local/include/googleapis \
    					--go_out=/app/generated/$$dir \
    					--go_opt=paths=source_relative \
    					--go-grpc_out=/app/generated/$$dir \
    					--go-grpc_opt=paths=source_relative \
    					$$file; \
    			done; \
    		done; \
    	'

## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
