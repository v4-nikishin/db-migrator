BIN:= "./bin/db-migrator"
GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./... 

test:
	go test -race ./internal/... -count=1 -v

integration-test:
	docker-compose -f docker-compose-testing.yaml up -d --build
	sleep 10 || echo "Waiting for prepear environment..."
	go test ./tests/... -count=1 -v || echo "Integration tests"
	docker-compose -f docker-compose-testing.yaml down

up:
	docker-compose up -d --build

down:
	docker-compose down
