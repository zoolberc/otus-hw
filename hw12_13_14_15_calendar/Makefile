BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"
DB_MIGRATION_CRED_POSTGRES="user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disabel"

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
	ls -a
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.54.2

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint

install-goose-deps:
	(which goose > /dev/null)|| curl -sSfL https://raw.githubusercontent.com/pressly/goose/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v3.5.0

migrate:install-goose-deps
	goose dir migrations postgres $(DB_MIGRATION_CRED_POSTGRES) up

migrate-status:install-goose-deps
	goose dir migrations postgres $(DB_MIGRATION_CRED_POSTGRES) status

migrate-down:install-goose-deps
	goose dir migrations postgres $(DB_MIGRATION_CRED_POSTGRES) down