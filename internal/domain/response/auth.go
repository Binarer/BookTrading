package response

// TokenResponse represents JWT token response
// @Description Ответ с токенами доступа
type TokenResponse struct {
	// @Description JWT токен доступа
	// @example eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	Token string `json:"token"`
	// @Description Токен для обновления доступа
	// @example eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	RefreshToken string `json:"refresh_token"`
}

// LoginResponse представляет полный ответ при входе пользователя
// @Description Полный ответ при успешном входе в систему
type LoginResponse struct {
	// @Description JWT токен доступа
	// @example eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	Token string `json:"token"`
	// @Description Токен для обновления доступа
	// @example eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	RefreshToken string `json:"refresh_token"`
	// @Description ID пользователя
	// @example 1
	UserID uint `json:"user_id"`
}
