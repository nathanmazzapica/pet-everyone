SOURCE := cmd/web/main.go
APP_NAME := pet-everyone
BIN_DIR := bin

run:
	go run ./...

pr-check:
	go fmt ./...
	go vet ./...
	staticcheck ./...



.PHONY: build clean linux

build:
	go build -o $(BIN_DIR)/$(APP_NAME)

linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	go build  -o $(BIN_DIR)/$(APP_NAME) $(SOURCE)

clean:
	rm -rf $(BIN_DIR)

