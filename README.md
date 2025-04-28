# Book Trading API

API для системы обмена книгами с поддержкой тегов.

## Технологии

- Go 1.21
- MySQL 8.0
- Chi Router
- Swagger
- Docker

## Установка

### Локальная установка
1. Клонируйте репозиторий
2. Установите зависимости:
```bash
go mod download
```
3. Создайте базу данных MySQL
4. Настройте переменные окружения (см. `.env.example`)
5. Запустите миграции:
```bash
go run cmd/migrate/main.go
```
6. Запустите приложение:
```bash
go run cmd/main.go
```

### Docker
1. Клонируйте репозиторий
2. Соберите и запустите контейнеры:
```bash
docker-compose up --build
```
Приложение будет доступно по адресу: http://localhost:8000

## API Endpoints

### Книги
- `POST /api/v1/books` - Создание книги
- `GET /api/v1/books/{id}` - Получение книги по ID
- `GET /api/v1/books/search` - Поиск книг по тегам
- `POST /api/v1/books/{id}/tags` - Добавление тегов к книге
- `PUT /api/v1/books/{id}` - Обновление книги
- `PATCH /api/v1/books/{id}/state` - Обновление состояния книги
- `DELETE /api/v1/books/{id}` - Удаление книги

### Теги
- `POST /api/v1/tags` - Создание тега
- `GET /api/v1/tags` - Получение всех тегов
- `GET /api/v1/tags/{id}` - Получение тега по ID
- `GET /api/v1/tags/popular` - Получение популярных тегов
- `DELETE /api/v1/tags/{id}` - Удаление тега

### Состояния
- `POST /api/v1/states` - Создание состояния
- `GET /api/v1/states` - Получение всех состояний
- `GET /api/v1/states/{id}` - Получение состояния по ID
- `PUT /api/v1/states/{id}` - Обновление состояния
- `DELETE /api/v1/states/{id}` - Удаление состояния

## Миграции
Миграции запускаются автоматически при старте Docker-контейнеров. Для ручного запуска:
```bash
go run cmd/migrate/main.go
```

## Swagger
Документация API доступна по адресу: http://localhost:8000/swagger/index.html

Для обновления документации:
```bash
swag init -g cmd/main.go -o docs
```

## Настройка Swagger
Для локальной разработки:
```bash
swag init -g cmd/main.go -o docs --host localhost:8000
```

Для продакшена:
```bash
swag init -g cmd/main.go -o docs --host your-domain.com
```

## Состояния книг

- `available` - книга доступна для обмена
- `trading` - книга находится в процессе обмена
- `traded` - книга обменяна

#### Создать состояние
```http
POST /api/v1/states
Content-Type: application/json

{
    "name": "available"
}
```

Ответ:
```json
{
    "id": 1,
    "name": "available",
    "created_at": "2025-04-28T12:00:00Z",
    "updated_at": "2025-04-28T12:00:00Z"
}
```

#### Получить все состояния
```http
GET /api/v1/states
```

Ответ:
```json
[
    {
        "id": 1,
        "name": "available",
        "created_at": "2025-04-28T12:00:00Z",
        "updated_at": "2025-04-28T12:00:00Z"
    },
    {
        "id": 2,
        "name": "trading",
        "created_at": "2025-04-28T12:00:00Z",
        "updated_at": "2025-04-28T12:00:00Z"
    }
]
```

#### Получить состояние по ID
```http
GET /api/v1/states/{id}
```

Ответ:
```json
{
    "id": 1,
    "name": "available",
    "created_at": "2025-04-28T12:00:00Z",
    "updated_at": "2025-04-28T12:00:00Z"
}
```

#### Обновить состояние
```http
PUT /api/v1/states/{id}
Content-Type: application/json

{
    "name": "traded"
}
```

Ответ:
```json
{
    "id": 1,
    "name": "traded",
    "created_at": "2025-04-28T12:00:00Z",
    "updated_at": "2025-04-28T12:00:00Z"
}
```

#### Удалить состояние
```http
DELETE /api/v1/states/{id}
```

#### Удалить книгу
```http
DELETE /api/v1/books/{id}
```

#### Удалить тег
```http
DELETE /api/v1/tags/{id}
```

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