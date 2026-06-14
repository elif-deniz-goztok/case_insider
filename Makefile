DB_NAME ?= case_insider

.PHONY: run build test fmt vet migrate seed reset-db

run:
	go run main.go

build:
	go build -o bin/case_insider .

test:
	go test ./...

fmt:
	gofmt -w .

vet:
	go vet ./...

migrate:
	psql $(DB_NAME) < db/schema.sql

seed:
	psql $(DB_NAME) < db/seed.sql

reset-db:
	psql $(DB_NAME) < db/schema.sql && psql $(DB_NAME) < db/seed.sql
