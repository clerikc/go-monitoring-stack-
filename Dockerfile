# Build stage
FROM golang:1.21-alpine AS builder

# Устанавливаем git (обязательно!)
RUN apk add --no-cache git

WORKDIR /app

# Копируем оба файла модуля
COPY src/go.mod src/go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем исходный код
COPY src/ ./

# Компилируем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /main .

# Run stage
FROM alpine:latest
WORKDIR /
COPY --from=builder /main /main
EXPOSE 8080
ENTRYPOINT ["/main"]