DC=docker compose --file ./docker/docker-compose.yml -p feature-flag-service

.PHONY: up
up:
	$(DC) up

.PHONY: down
down:
	$(DC) down

.PHONY: dv
dv:
	$(DC) down -v

.PHONY: purge
purge:
	$(DC) down --rmi=all --volumes --remove-orphans

new-migration:
	touch ./migrations/$$(date +%s)_$(name).sql

.PHONY: migrate
migrate:
	@go run cmd/migrate/main.go

.PHONY: cli
cli:
	@go run cmd/feature-flag-cli/main.go

.PHONY: api
api:
	@go run cmd/feature-flag-service-api/main.go