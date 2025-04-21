# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Установка зависимостей для сборки
RUN apk add --no-cache git gcc musl-dev

# Копируем только файлы модулей сначала для лучшего кэширования
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные исходники
COPY . .

# Запуск тестов
RUN go test -coverprofile=coverage.out ./...

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o pvz_service ./cmd/server

# Runtime stage
FROM alpine:3.18
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Копируем бинарник и миграции
COPY --from=builder /app/pvz_service .
COPY --from=builder /app/internal/storage/postgres/migrations ./migrations

# Настройки времени выполнения
EXPOSE 8080 3000 9000
CMD ["./pvz_service"]