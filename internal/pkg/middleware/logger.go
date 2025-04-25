package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// LoggerMiddleware логирует HTTP запросы
func LoggerMiddleware(next http.Handler) http.Handler {
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