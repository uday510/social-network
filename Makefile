include .envrc

MIGRATIONS_PATH := ./cmd/migrate/migrations
DB_MIGRATOR_ADDR ?= $(shell echo $$DB_MIGRATOR_ADDR)

.PHONY: migrate-create migrate-up migrate-down

# Create a new migration: make migrate-create name=create_users
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $$name

# Run up migrations
migrate-up:
	migrate -path=$(MIGRATIONS_PATH) -database="$(DB_MIGRATOR_ADDR)" up

# Run down migrations (you can specify how many, e.g., make migrate-down steps=1)
migrate-down:
	migrate -path=$(MIGRATIONS_PATH) -database="$(DB_MIGRATOR_ADDR)" down $(steps)


.PHONY: seed
seed:
	@go run cmd/migrate/seed/main.go