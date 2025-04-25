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

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs

EXPOSE 8000

CMD ["./main"] 