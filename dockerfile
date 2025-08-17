# Используем официальный образ Go
FROM golang:1.21-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение (новая точка входа)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app

# Используем минимальный образ для финального контейнера
FROM alpine:latest

# Устанавливаем ca-certificates для HTTPS запросов
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем собранное приложение
COPY --from=builder /app/main .

# Открываем порт
EXPOSE 8081

# Запускаем приложение
CMD ["./main"] 