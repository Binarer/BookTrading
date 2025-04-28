package main

import (
	"booktrading/internal/config"
	httpHandler "booktrading/internal/delivery/http"
	"booktrading/internal/domain/book"
	"booktrading/internal/domain/state"
	"booktrading/internal/domain/tag"
	"booktrading/internal/domain/user"
	"booktrading/internal/pkg/cache"
	"booktrading/internal/pkg/logger"
	"booktrading/internal/pkg/swagger"
	"booktrading/internal/repository/mysql"
	"booktrading/internal/usecase"
	"fmt"
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
	db, err := mysql.NewConnection(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", err)
	}

	// Автомиграция моделей
	if err := db.AutoMigrate(
		&book.Book{},
		&tag.Tag{},
		&state.State{},
		&user.User{},
	); err != nil {
		logger.Fatal("Failed to run migrations", err)
	}

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
	r.Use(loggerMiddleware)

	// Инициализация маршрутов
	handler.InitRoutes(r)

	// Запуск сервера
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}

	logger.Info(fmt.Sprintf("Server is running on port %d", cfg.Server.Port))
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal("Failed to start server", err)
	}
}

// loggerMiddleware логирует HTTP запросы
func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Создаем кастомный ResponseWriter для отслеживания статус-кода
		ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(ww, r)

		duration := time.Since(start)
		logger.Info(fmt.Sprintf("HTTP request: method=%s path=%s status=%d duration=%v",
			r.Method,
			r.URL.Path,
			ww.statusCode,
			duration,
		))
	})
}

// responseWriter кастомный ResponseWriter для отслеживания статус-кода
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
