.PHONY: build test run clean

BINARY_NAME=dangeon
CONFIG_PATH ?= config/config.json
EVENTS_PATH ?= events

build:
	go build -o $(BINARY_NAME) cmd/dangeon/main.go

test:
	go test ./...

cover:
	go test -coverprofile=./tests/coverage.out ./...
	go tool cover -func=./tests/coverage.out

run: build
	./$(BINARY_NAME) --config $(CONFIG_PATH) --events $(EVENTS_PATH)

clean:
	rm -f $(BINARY_NAME)
	go clean
