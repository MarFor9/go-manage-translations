LOCAL_DEV_PATH = $(shell pwd)/infrastructure/local

DOCKER_COMPOSE_FILE := $(LOCAL_DEV_PATH)/docker-compose.yml
DOCKER_COMPOSE_FILE_INFRA := $(LOCAL_DEV_PATH)/docker-compose-infra.yml

DOCKER_COMPOSE_CMD := docker compose -p translation -f $(DOCKER_COMPOSE_FILE)
DOCKER_COMPOSE_INFRA_CMD := docker compose -p translation -f $(DOCKER_COMPOSE_FILE_INFRA)

BIN := $(shell pwd)/bin
GO?=$(shell which go)
export GOBIN := $(BIN)
export PATH := $(BIN):$(PATH)

.PHONY: up
up:
	$(DOCKER_COMPOSE_INFRA_CMD) up -d postgres

.PHONY: up-test
up-test:
	$(DOCKER_COMPOSE_INFRA_CMD) up -d test_postgres

.PHONY: migrate
migrate: ## Run the database migration.
	$(DOCKER_COMPOSE_CMD) up -d migrate
	sleep 5
	docker logs translation-migrate-1
	docker rm translation-migrate-1

.PHONY: run
run: ## Run the docker image.
	$(DOCKER_COMPOSE_CMD) up -d api

.PHONY: down
down: ## Stop the docker image.
	$(DOCKER_COMPOSE_CMD) down --remove-orphans

.PHONY: tests
tests:
	$(GO) test -v ./... --count=1

.PHONY: api
api: ## Generate API files.
	oapi-codegen -o ./internal/api/api.gen.go -config api/config-oapi-codegen.yaml api/api.yaml