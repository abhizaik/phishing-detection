# === Makefile for Go + Svelte App ===
# Project structure:
# - Docker Compose files in ./docker/
# - Go backend in ./server/
# - Svelte frontend in ./web/website

# === Variables ===
DOCKER_DEV = docker/docker-compose.dev.yml
DOCKER_PROD = docker/docker-compose.prod.yml

# === Development Commands ===

## Start dev environment (backend via Docker Compose)
dev-up:
	docker compose -f $(DOCKER_DEV) up --build -d

## Stop dev environment
dev-down:
	docker compose -f $(DOCKER_DEV) down

## Rebuild containers and restart dev environment
dev-rebuild:
	docker compose -f $(DOCKER_DEV) up --build --force-recreate -d

## Tail logs for dev environment
dev-logs:
	docker compose -f $(DOCKER_DEV) logs -f

# === Production Commands ===

## Start prod environment (Docker Compose)
prod-up:
	docker compose -f $(DOCKER_PROD) up --build -d

## Stop prod environment
prod-down:
	docker compose -f $(DOCKER_PROD) down

## Rebuild and restart prod containers
prod-rebuild:
	docker compose -f $(DOCKER_PROD) up --build --force-recreate -d

## Tail logs for prod environment
prod-logs:
	docker compose -f $(DOCKER_PROD) logs -f

# === Local Development (Without Docker) ===

## Run Go backend locally
run-backend:
	cd server && air

## Run Svelte frontend locally (dev server)
run-frontend:
	cd web/website && npm run dev

## Build Svelte frontend
build-frontend:
	cd web/website && npm install && npm run build

# === Backend Testing & Linting ===

## Run Go tests
test-backend:
	cd server && go test ./...

## Run Go linter
lint-backend:
	cd server && go vet ./...

# === Cleanup ===

## Stop and remove all dev/prod containers and volumes, then prune unused Docker data
clean:
	docker-compose -f $(DOCKER_DEV) down -v
	docker-compose -f $(DOCKER_PROD) down -v
	docker system prune -f

# === Help ===

## Show available Make commands
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
