include .env

MIGRATIONS_PATH=cmd/migrate

migrate:
	@make migrate-down
	@make migrate-up

migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) up

migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) down $(filter-out $@,$(MAKECMDGOALS))

build:
	@go build -o bin/api ./cmd/api

run:build
	@./bin/api
