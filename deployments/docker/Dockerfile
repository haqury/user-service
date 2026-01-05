FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем файлы модулей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /user-service ./cmd/user-service

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Копируем бинарник из builder
COPY --from=builder /user-service .
COPY config.yaml ./

# Экспортируем порты
EXPOSE 8081 9091

# Запускаем приложение
CMD ["./user-service", "--config", "config.yaml"]
