package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// NewRouter создает новый роутер с настроенными маршрутами
func NewRouter(h *Handler, jwtAuth *jwtauth.JWTAuth) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.URLFormat)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	// Кастомный CORS middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Получаем origin из запроса
			origin := r.Header.Get("Origin")

			// Если origin есть в запросе, используем его
			if origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
			} else {
				// Если origin нет, разрешаем все origins
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}

			// Устанавливаем остальные CORS заголовки
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "300")

			// Если это preflight запрос, сразу отвечаем OK
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/api/v1/books", h.getAllBooks)
		r.Get("/api/v1/books/{id}", h.getBookByID)
		r.Get("/api/v1/books/tags", h.searchBooksByTags)
		r.Get("/api/v1/tags", h.getAllTags)
		r.Get("/api/v1/tags/{id}", h.getTagByID)
		r.Get("/api/v1/states", h.getAllStates)
		r.Get("/api/v1/states/{id}", h.getStateByID)
		r.Post("/api/v1/users/register", h.Register)
		r.Post("/api/v1/users/login", h.Login)
		r.Post("/api/v1/auth/refresh", h.RefreshToken)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(jwtAuth))
		r.Use(jwtauth.Authenticator(jwtAuth))

		// Book routes
		r.Post("/api/v1/books", h.createBook)
		r.Put("/api/v1/books/{id}", h.updateBook)
		r.Delete("/api/v1/books/{id}", h.deleteBook)
		r.Put("/api/v1/books/{id}/state", h.updateBookState)

		// Tag routes
		r.Post("/api/v1/tags", h.createTag)
		r.Put("/api/v1/tags/{id}", h.updateTag)
		r.Delete("/api/v1/tags/{id}", h.deleteTag)

		// State routes
		r.Post("/api/v1/states", h.createState)
		r.Put("/api/v1/states/{id}", h.updateState)
		r.Delete("/api/v1/states/{id}", h.deleteState)

		// User routes
		r.Get("/api/v1/users", h.getAllUsers)
		r.Get("/api/v1/users/{id}", h.GetUserByID)
		r.Put("/api/v1/users/{id}", h.UpdateUser)
		r.Delete("/api/v1/users/{id}", h.DeleteUser)
		r.Get("/api/v1/users/{id}/books", h.getUserBooks)
	})

	// Swagger
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	return r
}
