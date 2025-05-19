package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (h *Handler) InitRouter() *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	// Swagger
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Group(func(r chi.Router) {
			// User routes
			r.Post("/users/register", h.registerUser)
			r.Post("/users/login", h.loginUser)

			// Book routes
			r.Get("/books", h.getAllBooks)
			r.Get("/books/{id}", h.getBookByID)
			r.Get("/books/search", h.searchBooksByTags)

			// Tag routes
			r.Get("/tags", h.getAllTags)
			r.Get("/tags/{id}", h.getTagByID)
			r.Get("/tags/popular", h.getPopularTags)

			// State routes
			r.Get("/states", h.getAllStates)
			r.Get("/states/{id}", h.getStateByID)
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(h.AuthMiddleware)

			// User routes
			r.Get("/users", h.getAllUsers)
			r.Get("/users/{id}", h.getUserByID)
			r.Put("/users/{id}", h.updateUser)
			r.Delete("/users/{id}", h.deleteUser)

			// Book routes
			r.Post("/books", h.createBook)
			r.Put("/books/{id}", h.updateBook)
			r.Delete("/books/{id}", h.deleteBook)
			r.Post("/books/{id}/tags", h.addTagsToBook)
			r.Patch("/books/{id}/state", h.updateBookState)

			// Tag routes
			r.Post("/tags", h.createTag)
			r.Delete("/tags/{id}", h.deleteTag)

			// State routes
			r.Post("/states", h.createState)
			r.Put("/states/{id}", h.updateState)
			r.Delete("/states/{id}", h.deleteState)
		})
	})

	return r
}
