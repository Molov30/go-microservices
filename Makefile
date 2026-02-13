PROJECT_ROOT := $(abspath .)
GENERATED_DIR := $(PROJECT_ROOT)/generated
PROTO_DIR := $(PROJECT_ROOT)/proto

.PHONY: gen help

## gen: Generate Go code from proto files
gen:
	$(MAKE) -C $(PROTO_DIR) gen GENERATED_DIR=$(GENERATED_DIR) PROJECT_ROOT=$(PROJECT_ROOT)

## help: Print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
