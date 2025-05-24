package main

import (
	"booktrading/internal/config"
	httpHandler "booktrading/internal/delivery/http"
	"booktrading/internal/pkg/cache"
	"booktrading/internal/pkg/jwt"
	"booktrading/internal/pkg/logger"
	"booktrading/internal/repository"
	"booktrading/internal/repository/mysql"
	"booktrading/internal/usecase"
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "booktrading/docs" // This is required for Swagger
)

// @title Book Trading API
// @version 1.0
// @description API для обмена книгами
// @host localhost:8000
// @BasePath /
// @schemes http

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @securityDefinitions.apikey RefreshToken
// @in header
// @name X-Refresh-Token
// @description Refresh token for getting new access token.

func main() {
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load config", err)
	}

	// Получаем *gorm.DB напрямую
	db, err := mysql.InitDB(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	))
	if err != nil {
		logger.Fatal("Failed to connect to database", err)
	}

	repo := repository.NewRepository(db)

	jwtSvc := jwt.NewService(
		cfg.JWT.SecretKey,
		cfg.JWT.RefreshSecret,
		cfg.JWT.Issuer,
		repo.Token,
		repo.User,
	)

	// Инициализация кеша
	cache := cache.NewCache()

	// Инициализация usecase'ов
	bookUsecase := usecase.NewBookUsecase(
		repo.Book.(*mysql.BookRepository),
		repo.Tag.(*mysql.TagRepository),
		repo.State.(*mysql.StateRepository),
		cache,
	)

	tagUsecase := usecase.NewTagUsecase(
		repo.Tag.(*mysql.TagRepository),
		repo.Book.(*mysql.BookRepository),
		cache,
	)

	stateUsecase := usecase.NewStateUsecase(repo.State.(*mysql.StateRepository))

	userUsecase := usecase.NewUserUsecase(repo.User, jwtSvc)

	// Инициализация HTTP обработчика
	handler := httpHandler.NewHandler(
		bookUsecase,
		tagUsecase,
		stateUsecase,
		userUsecase,
		jwtSvc,
	)

	// Инициализация роутера
	router := httpHandler.NewRouter(handler, jwtSvc.GetJWTAuth())

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
