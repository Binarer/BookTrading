FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

# Установка swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

# Генерация Swagger документации с правильным хостом
ENV SWAGGER_HOST=10.3.13.28:8000
RUN swag init -g cmd/main.go -o docs

# Сборка приложения и миграций
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate ./cmd/migrate/main.go

FROM alpine:latest

WORKDIR /app

# Копируем бинарные файлы
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/migrations ./migrations

# Создаем скрипт для запуска
RUN echo '#!/bin/sh\n\
# Ждем, пока MySQL будет готов\n\
echo "Waiting for MySQL..."\n\
while ! nc -z mysql 3306; do\n\
  sleep 1\n\
done\n\
echo "MySQL is ready!"\n\
\n\
# Применяем миграции\n\
echo "Running migrations..."\n\
./migrate\n\
\n\
# Запускаем приложение\n\
echo "Starting application..."\n\
./main\n\
' > /app/start.sh && chmod +x /app/start.sh

# Устанавливаем netcat для проверки доступности MySQL
RUN apk add --no-cache netcat-openbsd

EXPOSE 8000

CMD ["/app/start.sh"] 