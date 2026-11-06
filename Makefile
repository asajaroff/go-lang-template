.PHONY: changelog install-tools build run test lint fmt sqlc migrate-up migrate-down docker-build docker-run compose-up compose-down helm-lint helm-template

install-tools:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0

build:
	go build -o bin/server ./cmd/server

run:
	go run ./cmd/server serve

test:
	go test ./...

lint:
	go vet ./...
	golangci-lint run ./...

fmt:
	gofmt -l -w .
	golangci-lint fmt ./...

sqlc:
	sqlc generate

migrate-up:
	go run ./cmd/server migrate up

migrate-down:
	go run ./cmd/server migrate down

docker-build:
	docker build -t go-lang-template:local .

docker-run:
	docker run --rm -p 8000:8000 --env-file .env go-lang-template:local

compose-up:
	docker compose up --build

compose-down:
	docker compose down -v

helm-lint:
	helm lint charts/go-lang-template

helm-template:
	helm template go-lang-template charts/go-lang-template

# Usage: make changelog TAG=v0.1.0 (needs Docker)
changelog:
	@test -n "$(TAG)" || (echo "usage: make changelog TAG=vX.Y.Z" && exit 1)
	docker run --rm -v "$(CURDIR)":/app/ orhunp/git-cliff:latest --tag $(TAG) -o CHANGELOG.md
