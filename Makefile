.PHONY: build, run, docker-build, docker-run, up, down

# Build app
build:
	go build -o bin/order-service ./cmd/app

# Launch app
run:
	go run ./cmd/app

# Building Docker image
docker-build:
	docker build -t app .

# Launching Docker container
docker-run:
	docker run -p 8081:8081 --env-file .env app

# Launching PostgreSQL, Kafka
up:
	docker-compose up -d

# Stopping PostgreSQL, Kafka
down:
	docker-compose down

# Check docker status
status:
	docker-compose ps

# logs
logs:
	docker-compose logs -f


