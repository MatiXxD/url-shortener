include .dev.env

MIGRATIONS_PATH := ./migrations
POSTGRES_CONNECTION = postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}/${POSTGRES_DB}?sslmode=disable

.PHONY: build-server
build-server:
	go build -o bin/server ./cmd/server

.PHONY: build-client
build-client:
	go build -o bin/client ./cmd/client

.PHONY: test
test:
	go test ./...

.PHONY: print-dsn
print-dsn:
	@echo ${POSTGRES_CONNECTION}

.PHONY: postgres-create
postgres-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: 'name' is not set. Usage: make postgres-create name=migration_name"; \
		exit 1; \
	fi
	@goose create -dir $(MIGRATIONS_PATH) $(name) sql

.PHONY: postgres-up
postgres-up:
	@GOOSE_DRIVER=postgres \
	GOOSE_DBSTRING="$(POSTGRES_CONNECTION)" \
	goose -dir $(MIGRATIONS_PATH) up

.PHONY: postgres-down
postgres-down:
	@GOOSE_DRIVER=postgres \
	GOOSE_DBSTRING="$(POSTGRES_CONNECTION)" \
	goose -dir $(MIGRATIONS_PATH) down

.PHONY: postgres-down-all
postgres-down-all:
	@GOOSE_DRIVER=postgres \
	GOOSE_DBSTRING="$(POSTGRES_CONNECTION)" \
	goose -dir $(MIGRATIONS_PATH) down-to 0
