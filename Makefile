.PHONY: install-lint-deps lint tidy install-protoc install-protoc-gen-go install-protoc-gen-go-grpc install-gomock generate \
build-server dc-up dc-down-prune dc-down makemigrations integration-tests build-cli build run

DCF := -f deployments/docker-compose.yaml
DCF_TEST := -p test -f deployments/docker-compose.integration.yaml

SERVICE_BIN := "./bin/service"
CLI_BIN := "./bin/cli"

install-lint-deps:
	@(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.53.3

lint: install-lint-deps
	golangci-lint run ./...

tidy:
	go mod tidy

install-protoc:
	@(which protoc > /dev/null) || (curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v24.2/protoc-24.2-linux-x86_64.zip && \
    unzip -o protoc-24.2-linux-x86_64.zip -d ${HOME}/.local &&\
    rm -rf protoc-24.2-linux-x86_64.zip)

install-protoc-gen-go:
	@(which protoc-gen-go > /dev/null) || go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.30.0

install-protoc-gen-go-grpc:
	@(which protoc-gen-go-grpc > /dev/null) || go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

install-gomock:
	@(which mockgen > /dev/null) || go install go.uber.org/mock/mockgen@v0.2.0

generate: install-protoc install-protoc-gen-go install-protoc-gen-go-grpc install-gomock
	go generate ./...

build-server:
	go build -v -o $(SERVICE_BIN) ./cmd/service/main.go

build-cli:
	go build -v -o $(CLI_BIN) ./cmd/cli/main.go

build: build-server build-cli

run-server: build-server
	APP_PG_DSN=postgres://user:password@localhost:15432/db APP_REDIS_DSN=redis://localhost:16379 $(SERVICE_BIN)

makemigrations: ## e.g `make makemigrations name=init`
	docker compose $(DCF) run --rm --no-deps migrator migrate create -ext sql -dir /migrations $(name)

dc-up: dc-down
	docker compose $(DCF) up --build --detach

dc-down:
	docker compose $(DCF) down

dc-down-prune:
	docker compose $(DCF) down -v --rmi all

test:
	go test -v -race -count=100 ./internal/delivery/grpc ./internal/usecase

integration-tests:
	docker compose $(DCF_TEST) down -v
	docker compose $(DCF_TEST) run --build --rm test; EXIT_CODE=$$?; docker compose $(DCF_TEST) down -v --rmi local; exit $${EXIT_CODE}

run: dc-up