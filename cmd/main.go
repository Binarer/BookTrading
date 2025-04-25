package main

import (
	"booktrading/internal/pkg/cache"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/rs/zerolog/log"

	"booktrading/internal/config"
	httpHandler "booktrading/internal/delivery/http"

	"booktrading/internal/pkg/logger"
	"booktrading/internal/pkg/middleware"
	"booktrading/internal/repository/mysql"
	"booktrading/internal/usecase"
)

// @title Book Trading API
// @version 1.0
// @description API for book trading system with tag support
// @host 10.3.13.28:8080
// @BasePath /api/v1
func main() {
	// Инициализация конфигурации
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Инициализация логгера
	logger.InitLogger(cfg.Logging.Level, cfg.Logging.Format)

	// Инициализация кэша
	localCache := cache.NewCache(cfg.Cache.TTL, cfg.Cache.CleanupInterval)

	// Подключение к базе данных
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	// Инициализация репозиториев
	tagRepo := mysql.NewTagRepository(db.DB)
	bookRepo := mysql.NewBookRepository(db.DB)

	// Инициализация use cases
	tagUsecase := usecase.NewTagUsecase(tagRepo, localCache)
	bookUsecase := usecase.NewBookUsecase(bookRepo, tagRepo, localCache)

	// Инициализация HTTP обработчика
	handler := httpHandler.NewHandler(tagUsecase, bookUsecase)

	// Настройка CORS
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{cfg.CORS.AllowedOrigins},
		AllowedMethods:   []string{cfg.CORS.AllowedMethods},
		AllowedHeaders:   []string{cfg.CORS.AllowedHeaders},
		ExposedHeaders:   []string{cfg.CORS.ExposedHeaders},
		AllowCredentials: cfg.CORS.AllowCredentials,
		MaxAge:           cfg.CORS.MaxAge,
	})

	// Создание роутера
	r := chi.NewRouter()
	r.Use(corsMiddleware.Handler)
	r.Use(middleware.LoggerMiddleware)

	// Инициализация маршрутов
	handler.InitRoutes(r)

	// Создание HTTP сервера
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Запуск сервера в горутине
	go func() {
		log.Info().Msgf("Starting server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	// Ожидание сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	log.Info().Msg("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exiting")
}

// loggerMiddleware логирует HTTP запросы
func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Создаем кастомный ResponseWriter для отслеживания статус-кода
		ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(ww, r)

		duration := time.Since(start)
		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", ww.statusCode).
			Dur("duration", duration).
			Msg("HTTP request")
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
