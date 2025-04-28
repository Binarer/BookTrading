# Book Trading API

API для системы обмена книгами с поддержкой тегов.

## Технологии

- Go 1.21
- MySQL 8.0
- Docker
- Swagger

## Установка и запуск

### С помощью Docker

1. Клонируйте репозиторий:
```bash
git clone https://github.com/yourusername/booktrading.git
cd booktrading
```

2. Запустите приложение с помощью Docker Compose:
```bash
docker-compose up -d
```

Приложение будет доступно по адресу: http://localhost:8000

### Локальная установка

1. Установите зависимости:
```bash
go mod download
```

2. Создайте базу данных:
```bash
mysql -u root -p < migrations/001_initial_schema.sql
```

3. Запустите приложение:
```bash
go run cmd/main.go
```

## API Endpoints

### Теги

#### Создать тег
```http
POST /api/v1/tags
Content-Type: application/json

{
    "name": "fiction"
}
```

Ответ:
```json
{
    "id": 1,
    "name": "fiction",
    "created_at": "2025-04-28T12:00:00Z",
    "updated_at": "2025-04-28T12:00:00Z"
}
```

#### Получить тег по ID
```http
GET /api/v1/tags/{id}
```

Ответ:
```json
{
    "id": 1,
    "name": "fiction",
    "created_at": "2025-04-28T12:00:00Z",
    "updated_at": "2025-04-28T12:00:00Z"
}
```

#### Получить популярные теги
```http
GET /api/v1/tags/popular?limit=10
```

Ответ:
```json
[
    {
        "id": 1,
        "name": "fiction",
        "created_at": "2025-04-28T12:00:00Z",
        "updated_at": "2025-04-28T12:00:00Z"
    }
]
```

### Книги

#### Создать книгу
```http
POST /api/v1/books
Content-Type: application/json

{
    "title": "The Great Gatsby",
    "author": "F. Scott Fitzgerald",
    "description": "A story of the fabulously wealthy Jay Gatsby and his love for the beautiful Daisy Buchanan.",
    "state": "available",
    "photos": ["data:image/jpeg;base64,/9j/4AAQSkZJRg..."],
    "tag_ids": [1, 2]
}
```

Ответ:
```json
{
    "id": 1,
    "title": "The Great Gatsby",
    "author": "F. Scott Fitzgerald",
    "description": "A story of the fabulously wealthy Jay Gatsby and his love for the beautiful Daisy Buchanan.",
    "state": "available",
    "photos": ["data:image/jpeg;base64,/9j/4AAQSkZJRg..."],
    "created_at": "2025-04-28T12:00:00Z",
    "updated_at": "2025-04-28T12:00:00Z",
    "tags": [
        {
            "id": 1,
            "name": "fiction",
            "created_at": "2025-04-28T12:00:00Z",
            "updated_at": "2025-04-28T12:00:00Z"
        }
    ]
}
```

#### Получить книгу по ID
```http
GET /api/v1/books/{id}
```

Ответ:
```json
{
    "id": 1,
    "title": "The Great Gatsby",
    "author": "F. Scott Fitzgerald",
    "description": "A story of the fabulously wealthy Jay Gatsby and his love for the beautiful Daisy Buchanan.",
    "state": "available",
    "photos": ["data:image/jpeg;base64,/9j/4AAQSkZJRg..."],
    "created_at": "2025-04-28T12:00:00Z",
    "updated_at": "2025-04-28T12:00:00Z",
    "tags": [
        {
            "id": 1,
            "name": "fiction",
            "created_at": "2025-04-28T12:00:00Z",
            "updated_at": "2025-04-28T12:00:00Z"
        }
    ]
}
```

#### Поиск книг по тегам
```http
GET /api/v1/books/search?tag_id=1&tag_id=2
```

Ответ:
```json
[
    {
        "id": 1,
        "title": "The Great Gatsby",
        "author": "F. Scott Fitzgerald",
        "description": "A story of the fabulously wealthy Jay Gatsby and his love for the beautiful Daisy Buchanan.",
        "state": "available",
        "photos": ["data:image/jpeg;base64,/9j/4AAQSkZJRg..."],
        "created_at": "2025-04-28T12:00:00Z",
        "updated_at": "2025-04-28T12:00:00Z",
        "tags": [
            {
                "id": 1,
                "name": "fiction",
                "created_at": "2025-04-28T12:00:00Z",
                "updated_at": "2025-04-28T12:00:00Z"
            }
        ]
    }
]
```

#### Добавить теги к книге
```http
POST /api/v1/books/{id}/tags
Content-Type: application/json

[1, 2]
```

#### Обновить книгу
```http
PUT /api/v1/books/{id}
Content-Type: application/json

{
    "title": "The Great Gatsby",
    "author": "F. Scott Fitzgerald",
    "description": "A story of the fabulously wealthy Jay Gatsby and his love for the beautiful Daisy Buchanan.",
    "state": "trading",
    "photos": ["data:image/jpeg;base64,/9j/4AAQSkZJRg..."]
}
```

Ответ:
```json
{
    "id": 1,
    "title": "The Great Gatsby",
    "author": "F. Scott Fitzgerald",
    "description": "A story of the fabulously wealthy Jay Gatsby and his love for the beautiful Daisy Buchanan.",
    "state": "trading",
    "photos": ["data:image/jpeg;base64,/9j/4AAQSkZJRg..."],
    "created_at": "2025-04-28T12:00:00Z",
    "updated_at": "2025-04-28T12:00:00Z",
    "tags": [
        {
            "id": 1,
            "name": "fiction",
            "created_at": "2025-04-28T12:00:00Z",
            "updated_at": "2025-04-28T12:00:00Z"
        }
    ]
}
```

## Состояния книг

- `available` - книга доступна для обмена
- `trading` - книга находится в процессе обмена
- `traded` - книга обменяна

## Миграции

### Создание базы данных

1. Подключитесь к MySQL:
```bash
mysql -u root -p
```

2. Создайте базу данных:
```sql
CREATE DATABASE booktrading;
USE booktrading;
```

3. Примените миграции:
```bash
mysql -u root -p booktrading < migrations/001_initial_schema.sql
```

## Docker

### docker-compose.yml

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8000:8000"
    environment:
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USER=root
      - DB_PASSWORD=root
      - DB_NAME=booktrading
    depends_on:
      - mysql

  mysql:
    image: mysql:8.0
    ports:
      - "3306:3306"
    environment:
      - MYSQL_ROOT_PASSWORD=root
      - MYSQL_DATABASE=booktrading
    volumes:
      - mysql_data:/var/lib/mysql
      - ./migrations:/docker-entrypoint-initdb.d

volumes:
  mysql_data:
```

### Dockerfile

```dockerfile
FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

EXPOSE 8000

CMD ["./main"]
```

## Swagger

Документация API доступна по адресу: http://localhost:8000/swagger/index.html

Для обновления Swagger документации:
```bash
swag init -g cmd/main.go
```

### Конфигурация хоста Swagger

Для локальной разработки добавьте в файлы:
- `cmd/main.go`
- `internal/delivery/http/handler.go`

```go
// @host localhost:8000
```

Для продакшена (VM) добавьте:
```go
// @host 10.3.13.28:8000
```

После изменения хоста необходимо перегенерировать документацию:
```bash
swag init -g cmd/main.go
``` 