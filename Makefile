.PHONY: run build test test-integration test-all lint fmt vet vuln \
        migrate-up migrate-down migrate-create sqlc \
        docker-build docker-up docker-down \
        generate install-tools

# Go parameters
BINARY_NAME=server
MAIN_PATH=./cmd/server
COVERAGE_FILE=coverage.out

# Database
DB_URL ?= postgres://postgres:postgres@localhost:5432/boilerplate?sslmode=disable

install-tools:
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	go install github.com/google/wire/cmd/wire@latest

generate:
	oapi-codegen -config .oapi-codegen.yaml api/openapi.yaml
	wire ./internal/di

run:
	go run $(MAIN_PATH)

build:
	CGO_ENABLED=0 go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

test:
	go test -race -count=1 ./...

test-integration:
	go test -tags=integration -race -count=1 ./...

test-all:
	go test -tags=integration -race -count=1 -coverprofile=$(COVERAGE_FILE) ./...

coverage: test-all
	go tool cover -html=$(COVERAGE_FILE) -o coverage.html

lint:
	golangci-lint run ./...

fmt:
	gofumpt -w .

vet:
	go vet ./...

vuln:
	govulncheck ./...

tidy:
	go mod tidy

migrate-up:
	migrate -path internal/adapter/postgresrepo/migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path internal/adapter/postgresrepo/migrations -database "$(DB_URL)" down 1

migrate-create:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir internal/adapter/postgresrepo/migrations -seq $$name

sqlc:
	sqlc generate

docker-build:
	docker build -f deployments/Dockerfile -t golang-vibe-boilerplate:latest .

docker-up:
	docker compose -f deployments/docker-compose.yml up -d

docker-down:
	docker compose -f deployments/docker-compose.yml down
