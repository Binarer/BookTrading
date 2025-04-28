package main

import (
	"booktrading/internal/config"
	httpHandler "booktrading/internal/delivery/http"
	"booktrading/internal/pkg/cache"
	"booktrading/internal/pkg/logger"
	"booktrading/internal/repository/mysql"
	"booktrading/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

// @title Book Trading API
// @version 1.0
// @description API for book trading system with tag support
// @host localhost:8000
// @BasePath /
func main() {
	// Инициализация логгера
	logger.Init()

	// Загрузка конфигурации
	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatal("Failed to load config", err)
	}

	// Инициализация кеша
	cache := cache.NewCache(5*time.Minute, 10*time.Minute)

	// Инициализация репозиториев
	db, err := mysql.NewMySQLConnection(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", err)
	}
	defer db.Close()

	bookRepo := mysql.NewBookRepository(db)
	tagRepo := mysql.NewTagRepository(db)
	stateRepo := mysql.NewStateRepository(db)

	// Инициализация usecases
	bookUsecase := usecase.NewBookUsecase(bookRepo, tagRepo, cache)
	tagUsecase := usecase.NewTagUsecase(tagRepo, bookRepo, cache)
	stateUsecase := usecase.NewStateUsecase(stateRepo)

	// Инициализация HTTP обработчика
	handler := httpHandler.NewHandler(tagUsecase, bookUsecase, stateUsecase)

	// Инициализация роутера
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CleanPath)
	r.Use(middleware.GetHead)

	// Инициализация маршрутов
	handler.InitRoutes(r)

	// Запуск сервера
	logger.Info("Server starting on port 8000")
	if err := http.ListenAndServe(":8000", r); err != nil {
		logger.Fatal("Failed to start server", err)
	}
}
