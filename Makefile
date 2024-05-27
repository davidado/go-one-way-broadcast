MAIN_PACKAGE_PATH := ./cmd/redis
BINARY_NAME := broadcast

build:
	@go build -o bin/$(BINARY_NAME) $(MAIN_PACKAGE_PATH)/main.go

test:
	@go test -v ./...

run: build
	@./bin/$(BINARY_NAME)