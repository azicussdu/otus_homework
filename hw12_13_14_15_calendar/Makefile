BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

DB_URL := "postgres://azicussdu:azicussdu@localhost:5432/calendar_db?sslmode=disable"
MIGRATE_DIR="./migrations"

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.60.1

lint: install-lint-deps
	golangci-lint run ./...

migrate-up:
	goose -dir $(MIGRATE_DIR) postgres "$(DB_URL)" up

migrate-down:
	goose -dir $(MIGRATE_DIR) postgres "$(DB_URL)" down

migrate-status:
	goose -dir $(MIGRATE_DIR) postgres "$(DB_URL)" status

# command sample: make migrate-create name=alter_events_table
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: Please specify a migration name. Usage: make migrate-create name=example"; \
		exit 1; \
	fi
	goose -dir $(MIGRATE_DIR) create $(name) sql

clean:
	rm -f $(BIN)

help:
	@echo "  make build     - Compile the service binary"
	@echo "  make run       - Compile and run the service with the default config"
	@echo "  make test      - Run unit tests with the race detector"
	@echo "  make lint      - Run golangci-lint"
	@echo "  make migrate   - Apply database migrations (if applicable)"
	@echo "  make clean     - Remove build artifacts"
	@echo "  make help      - Display this help message"

.PHONY: build run build-img run-img version test lint migrate-up migrate-down migrate-status migrate-create clean help
