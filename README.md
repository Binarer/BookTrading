# Book Trading API

Book Trading API - это RESTful API для системы обмена книгами с поддержкой тегов. API позволяет пользователям создавать книги, добавлять к ним теги, искать книги по тегам и управлять популярными тегами.

## 🚀 Технологии

- **Go** - основной язык программирования
- **MySQL** - база данных
- **Chi** - легковесный роутер для Go
- **Swagger** - документация API
- **Cache** - кеширование для оптимизации производительности
- **Docker** - контейнеризация приложения

## 🚀 Запуск проекта

### Локальная установка

1. Клонируйте репозиторий:
```bash
git clone https://github.com/yourusername/booktrading.git
cd booktrading
```

2. Создайте файл `.env` в корне проекта:
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=booktrading
```

3. Запустите миграции:
```bash
go run cmd/migrate/main.go
```

4. Запустите сервер:
```bash
go run cmd/api/main.go
```

### Docker Compose

1. Клонируйте репозиторий:
```bash
git clone https://github.com/yourusername/booktrading.git
cd booktrading
```

2. Создайте файл `.env` в корне проекта:
```env
DB_HOST=mysql
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=booktrading
```

3. Запустите приложение:
```bash
docker-compose up --build
```

4. Для остановки:
```bash
docker-compose down
```

Сервер запустится на `http://localhost:8080`

### Docker Compose файл

```yaml
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USER=root
      - DB_PASSWORD=root
      - DB_NAME=booktrading
    depends_on:
      - mysql
    networks:
      - booktrading-network

  mysql:
    image: mysql:5.7
    ports:
      - "3306:3306"
    environment:
      - MYSQL_ROOT_PASSWORD=root
      - MYSQL_DATABASE=booktrading
    volumes:
      - mysql-data:/var/lib/mysql
    networks:
      - booktrading-network

networks:
  booktrading-network:
    driver: bridge

volumes:
  mysql-data:
```

## 📚 API Endpoints

### Теги

#### Создание тега
```http
POST /api/v1/tags
Content-Type: application/json

{
    "name": "fiction"
}
```

#### Получение тега по ID
```http
GET /api/v1/tags/{id}
```

#### Получение популярных тегов
```http
GET /api/v1/tags/popular?limit=10
```

### Книги

#### Создание книги
```http
POST /api/v1/books
Content-Type: application/json

{
    "title": "The Great Gatsby",
    "author": "F. Scott Fitzgerald",
    "description": "A story of the fabulously wealthy Jay Gatsby and his love for the beautiful Daisy Buchanan."
}
```

#### Получение книги по ID
```http
GET /api/v1/books/{id}
```

#### Поиск книг по тегам
```http
GET /api/v1/books/search?tag_id=1&tag_id=2
```

#### Добавление тегов к книге
```http
POST /api/v1/books/{id}/tags
Content-Type: application/json

{
    "tag_ids": [1, 2, 3]
}
```

## 📖 Документация API

Документация API доступна через Swagger UI:
```http
GET /swagger/
```

## 🛠 Разработка

### Структура проекта

```
booktrading/
├── cmd/                    # Точки входа
├── internal/              # Внутренний код
│   ├── domain/           # Доменные модели
│   ├── repository/       # Репозитории
│   ├── usecase/         # Бизнес-логика
│   ├── delivery/        # Доставка (HTTP, gRPC)
│   └── pkg/             # Вспомогательные пакеты
├── migrations/           # SQL миграции
└── docs/                # Документация
```

### Тестирование

```bash
go test ./...
```

## 📝 Лицензия

MIT License

## 🤝 Вклад в проект

1. Fork репозитория
2. Создайте ветку для вашей фичи (`git checkout -b feature/amazing-feature`)
3. Сделайте коммит ваших изменений (`git commit -m 'Add some amazing feature'`)
4. Запушьте в ветку (`git push origin feature/amazing-feature`)
5. Откройте Pull Request 