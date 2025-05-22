package main

import (
	"booktrading/internal/config"
	httpHandler "booktrading/internal/delivery/http"
	"booktrading/internal/pkg/cache"
	"booktrading/internal/pkg/jwt"
	"booktrading/internal/pkg/logger"
	"booktrading/internal/repository/mysql"
	"booktrading/internal/usecase"
	"fmt"
	"net/http"
	"strconv"
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
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load config", err)
	}

	// Формирование DSN для подключения к MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)

	// Инициализация базы данных с автоматической миграцией
	db, err := mysql.InitDB(dsn)
	if err != nil {
		logger.Fatal("Failed to initialize database", err)
	}

	// Инициализация кэша
	cacheInstance := cache.NewCache(cfg.Cache.TTL, cfg.Cache.CleanupInterval)

	// Инициализация JWT сервиса
	jwtSvc := jwt.NewService(cfg.JWT.SecretKey)

	// Инициализация репозиториев
	userRepo := mysql.NewUserRepository(db)
	bookRepo := mysql.NewBookRepository(db)
	tagRepo := mysql.NewTagRepository(db)
	stateRepo := mysql.NewStateRepository(db)

	// Инициализация use cases
	userUsecase := usecase.NewUserUsecase(userRepo, jwtSvc)
	bookUsecase := usecase.NewBookUsecase(bookRepo, tagRepo, cacheInstance)
	tagUsecase := usecase.NewTagUsecase(tagRepo, bookRepo, cacheInstance)
	stateUsecase := usecase.NewStateUsecase(stateRepo)

	// Инициализация HTTP обработчика
	handler := httpHandler.NewHandler(
		bookUsecase,
		tagUsecase,
		stateUsecase,
		userUsecase,
		jwtSvc,
	)

	// Инициализация роутера
	router := handler.InitRouter()

	// Запуск сервера
	addr := cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port)
	logger.Info("Server starting on " + addr)
	if err := http.ListenAndServe(addr, router); err != nil {
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
