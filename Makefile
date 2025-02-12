phony: migrate-add migrate-up migrate-down install-dependencies

migrate-add:
	@read -p "Enter the name of the migration: " name; \
		go run ./cmd/migrate create $$name sql

migrate-up:
	@go run ./cmd/migrate up

migrate-down:
	@go run ./cmd/migrate down

install-dependencies:
	@go install github.com/swaggo/swag/cmd/swag@latest

generate-swag:
	@swag init -o ./docs -d ./cmd/server/,./internal/