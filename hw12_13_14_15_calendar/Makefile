BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

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

migrate:
	goose -dir migrations postgres "host=localhost port=5432 user=azicussdu password=azicussdu dbname=calendar_db sslmode=disable" up

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

.PHONY: build run build-img run-img version test lint migrate clean help
