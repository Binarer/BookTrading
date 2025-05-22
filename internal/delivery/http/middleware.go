package http

import (
	"booktrading/internal/pkg/jwt"
	"booktrading/internal/pkg/logger"
	"context"
	"fmt"
	"net/http"
	"strings"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	LoginKey  contextKey = "login"
)

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Error("Authorization header is missing", nil)
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		logger.Info("Received Authorization header: " + authHeader)

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Error("Invalid authorization header format", nil)
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		logger.Info("Validating token: " + parts[1])

		claims, err := h.jwtSvc.ValidateToken(parts[1])
		if err != nil {
			if err == jwt.ErrExpiredToken {
				logger.Error("Token has expired", err)
				http.Error(w, "Token has expired", http.StatusUnauthorized)
				return
			}
			logger.Error("Invalid token", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		logger.Info(fmt.Sprintf("Token validated successfully for user ID: %d, login: %s", claims.UserID, claims.Login))

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, LoginKey, claims.Login)

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
