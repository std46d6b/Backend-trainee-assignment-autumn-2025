BINARY_NAME=app
CMD_PATH=./cmd/app

BIN_DIR=bin

.PHONY: build run test lint docker-up docker-down

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_PATH)

run:
	go run $(CMD_PATH)

test:
	go test ./...

lint:
	golangci-lint run

docker-up:
	docker compose -f docker-compose.yaml up --build -d

docker-down:
	docker compose -f docker-compose.yaml down

docker-logs:
	docker compose -f docker-compose.yaml logs

docker-logs-f:
	docker compose -f docker-compose.yaml logs -f

help:
	@echo "make build - build binary"
	@echo "make run - run binary"
	@echo "make test - run tests"
	@echo "make lint - run linter"
	@echo "make docker-up - start docker containers"
	@echo "make docker-down - stop docker containers"
	@echo "make docker-logs - show docker logs"
	@echo "make docker-logs-f - show docker logs in real time"
