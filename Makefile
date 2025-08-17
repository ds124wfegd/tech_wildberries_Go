.PHONY: build run test clean docker-build docker-run infra-up infra-down producer

# Сборка приложения
build:
	go build -o bin/order-service ./cmd/app

# Запуск приложения
run:
	go run ./cmd/app

# Тестирование
test:
	go test ./...

test-integration:
	go test -tags=integration ./internal/infra/postgres -v

# Очистка
clean:
	rm -rf bin/
	go clean

# Сборка Docker образа
docker-build:
	docker build -t app .

# Запуск Docker контейнера
docker-run:
	docker run -p 8081:8081 --env-file env.example app

# Запуск инфраструктуры (PostgreSQL, Kafka)
infra-up:
	docker-compose up -d

# Остановка инфраструктуры
infra-down:
	docker-compose down



# Проверка статуса
status:
	docker-compose ps

# Логи инфраструктуры
logs:
	docker-compose logs -f

