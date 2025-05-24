package http

import (
	"booktrading/internal/pkg/jwt"
	"booktrading/internal/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем заголовок Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Error("Authorization header is missing", nil)
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Проверяем формат заголовка
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Error("Invalid authorization header format", fmt.Errorf("expected 'Bearer <token>', got '%s'", authHeader))
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		if token == "" {
			logger.Error("Empty token", nil)
			http.Error(w, "Token is required", http.StatusUnauthorized)
			return
		}

		// Валидируем токен
		claims, err := h.jwtSvc.ValidateToken(token)
		if err != nil {
			switch err {
			case jwt.ErrExpiredToken:
				logger.Error("Token has expired", err)
				http.Error(w, "Token has expired", http.StatusUnauthorized)
			case jwt.ErrInvalidToken:
				logger.Error("Invalid token", err)
				http.Error(w, "Invalid token", http.StatusUnauthorized)
			default:
				logger.Error("Token validation failed", err)
				http.Error(w, "Authentication failed", http.StatusUnauthorized)
			}
			return
		}

		// Проверяем наличие обязательных полей в claims
		if claims.UserID == 0 {
			logger.Error("Invalid token claims", fmt.Errorf("user ID is missing or zero"))
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		if claims.Login == "" {
			logger.Error("Invalid token claims", fmt.Errorf("login is missing or empty"))
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Добавляем данные пользователя в контекст
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, LoginKey, claims.Login)

		// Добавляем заголовки безопасности
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromContext(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(UserIDKey).(uint)
	return userID, ok
}

func GetLoginFromContext(ctx context.Context) (string, bool) {
	login, ok := ctx.Value(LoginKey).(string)
	return login, ok
}
