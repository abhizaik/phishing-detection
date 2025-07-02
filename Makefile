APP_NAME := safesurf

up:
	docker compose -f docker/docker-compose.dev.yml up --build

down:
	docker compose -f docker/docker-compose.dev.yml down

build:
	docker build -f docker/Dockerfile.dev -t $(APP_NAME) .

run:
	docker run -p 8080:8080 $(APP_NAME)

test:
	go test ./...

tidy:
	go mod tidy
