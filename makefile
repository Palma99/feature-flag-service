COMPOSE=docker compose --file ./docker/docker-compose.yml $(COMPOSE_OVERRIDE) -p feature-flag-service

.PHONY: up
up:
	$(COMPOSE) up $$ARG

.PHONY: down
down:
	$(COMPOSE) down $$ARG

.PHONY: dv
dv:
	$(COMPOSE) down -v

.PHONY: purge
purge:
	$(COMPOSE) down --rmi=all --volumes --remove-orphans

new-migration:
	touch ./migrations/$$(date +%s)_$(name).sql

.PHONY: migrate
migrate:
	@go run cmd/migrate/main.go

.PHONY: cli
cli:
	@go run cmd/feature-flag-cli/main.go