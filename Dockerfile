FROM golang:1.21-alpine AS builder

WORKDIR /app

# Установка необходимых зависимостей
RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./

RUN go mod download

# Установка swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

# Генерация Swagger документации с правильным хостом
ENV SWAGGER_HOST=10.3.13.28:8000
RUN swag init -g cmd/main.go -o docs

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

FROM golang:1.21-alpine

WORKDIR /app

# Установка MySQL клиента
RUN apk add --no-cache mysql-client

COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/go.mod .
COPY --from=builder /app/go.sum .
COPY --from=builder /app/cmd ./cmd
COPY --from=builder /app/internal ./internal

EXPOSE 8000

CMD ["./main"] 