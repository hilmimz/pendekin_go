# Makefile for the Go project

APP_NAME=pendekin_go
BUILD_DIR=./bin
MAIN_PATH=./cmd/main.go

build:
	@go build -o $(BUILD_DIR)/$(APP_NAME).exe $(MAIN_PATH)

build-docker:
	@GOOS=linux GOARCH=amd64 go build -o ./bin/$(APP_NAME) ./cmd/main.go

run:
	@go run $(MAIN_PATH)