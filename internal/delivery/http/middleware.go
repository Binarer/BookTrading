package http

import (
	"booktrading/internal/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/jwtauth/v5"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	LoginKey  contextKey = "login"
)

// PhotoValidationMiddleware проверяет фотографии в запросе
func (h *Handler) PhotoValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем только POST и PUT запросы
		if r.Method != http.MethodPost && r.Method != http.MethodPut {
			next.ServeHTTP(w, r)
			return
		}

		// Проверяем только запросы к книгам
		if !strings.HasPrefix(r.URL.Path, "/api/books") {
			next.ServeHTTP(w, r)
			return
		}

		// Читаем тело запроса
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			logger.Error("Failed to decode request body", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Проверяем наличие фотографий
		photos, ok := body["photos"].([]interface{})
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		// Проверяем количество фотографий
		if len(photos) > 5 {
			logger.Error("Too many photos", fmt.Errorf("maximum 5 photos allowed, got %d", len(photos)))
			http.Error(w, "Maximum 5 photos allowed", http.StatusBadRequest)
			return
		}

		// Проверяем каждую фотографию
		for i, photo := range photos {
			photoStr, ok := photo.(string)
			if !ok {
				logger.Error("Invalid photo format", fmt.Errorf("photo at index %d is not a string", i))
				http.Error(w, fmt.Sprintf("Invalid photo format at index %d", i), http.StatusBadRequest)
				return
			}

			// Проверяем формат base64
			if !strings.HasPrefix(photoStr, "data:image/") {
				logger.Error("Invalid photo format", fmt.Errorf("photo at index %d is not a base64 image", i))
				http.Error(w, fmt.Sprintf("Invalid photo format at index %d: must be base64 encoded image", i), http.StatusBadRequest)
				return
			}

			// Проверяем формат изображения
			contentType := strings.TrimPrefix(photoStr, "data:")
			if !strings.HasPrefix(contentType, "image/jpeg;base64,") && !strings.HasPrefix(contentType, "image/png;base64,") {
				logger.Error("Unsupported image format", fmt.Errorf("photo at index %d has unsupported format", i))
				http.Error(w, fmt.Sprintf("Unsupported image format at index %d: only JPEG and PNG are allowed", i), http.StatusBadRequest)
				return
			}
		}

		// Восстанавливаем тело запроса
		newBody, err := json.Marshal(body)
		if err != nil {
			logger.Error("Failed to marshal request body", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Создаем новый запрос с восстановленным телом
		newReq, err := http.NewRequest(r.Method, r.URL.String(), strings.NewReader(string(newBody)))
		if err != nil {
			logger.Error("Failed to create new request", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Копируем заголовки
		for key, values := range r.Header {
			for _, value := range values {
				newReq.Header.Add(key, value)
			}
		}

		next.ServeHTTP(w, newReq)
	})
}

// GetUserIDFromContext извлекает ID пользователя из контекста запроса
func GetUserIDFromContext(ctx context.Context) (uint, bool) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil || claims == nil {
		return 0, false
	}

	userID, idOK := claims["user_id"].(float64)
	if !idOK {
		return 0, false
	}

	return uint(userID), true
}

// GetLoginFromContext извлекает логин пользователя из контекста запроса
func GetLoginFromContext(ctx context.Context) (string, bool) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil || claims == nil {
		return "", false
	}

	login, loginOK := claims["login"].(string)
	if !loginOK {
		return "", false
	}

	return login, true
}
