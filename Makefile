# Makefile for the Go project
include .env
export

APP_NAME=pendekin_go
BUILD_DIR=./bin
MAIN_PATH=./cmd/main.go
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

build:
	@go build -o $(BUILD_DIR)/$(APP_NAME).exe $(MAIN_PATH)

build-docker:
	@GOOS=linux GOARCH=amd64 go build -o ./bin/$(APP_NAME) ./cmd/main.go

run:
	@go run $(MAIN_PATH)

migrate-up:
	@migrate -path migrations -database "$(DB_URL)" up

migrate-force:
	@migrate -path migrations -database "$(DB_URL)" force $(version)

migrate-down:
	@migrate -path migrations -database "$(DB_URL)" down 1