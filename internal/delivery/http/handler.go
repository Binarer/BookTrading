package http

import (
	"booktrading/internal/domain/book"
	"booktrading/internal/domain/state"
	"booktrading/internal/domain/tag"
	"booktrading/internal/pkg/jwt"
	"booktrading/internal/pkg/logger"
	"booktrading/internal/pkg/validator"
	"booktrading/internal/usecase"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @x-codeSamples.languages ["curl", "python", "javascript"]

// @x-codeSamples.curl {"name": "Example with curl", "lang": "bash", "source": "curl -X POST \"http://localhost:8000/api/auth/login\" \\\n  -H \"Content-Type: application/json\" \\\n  -d '{\"login\":\"user\",\"password\":\"pass\"}'"}

// @x-codeSamples.python {"name": "Example with Python", "lang": "python", "source": "import requests\nresponse = requests.post(\"http://localhost:8000/api/auth/login\",\n  json={\"login\":\"user\",\"password\":\"pass\"})\nprint(response.json())"}

// @x-codeSamples.javascript {"name": "Example with JavaScript", "lang": "javascript", "source": "fetch(\"http://localhost:8000/api/auth/login\", {\n  method: \"POST\",\n  headers: {\"Content-Type\": \"application/json\"},\n  body: JSON.stringify({login:\"user\",password:\"pass\"})\n})\n.then(response => response.json())\n.then(data => console.log(data))"}

// @x-tryItOutEnabled true
// @x-validateRequest true
// @x-validateResponse true

// @tag.name Auth
// @tag.description Authentication operations

// @tag.name Users
// @tag.description User management operations

// @tag.name Books
// @tag.description Book management operations

// @tag.name Tags
// @tag.description Tag management operations

// @tag.name States
// @tag.description Book state management operations

// ErrorResponse представляет собой структуру для ответов с ошибками
// @Description Структура для возврата ошибок API
type ErrorResponse struct {
	Error string `json:"error" example:"Error message"`
}

// Handler представляет HTTP обработчик
type Handler struct {
	bookUsecase  usecase.BookUsecase
	tagUsecase   usecase.TagUsecase
	stateUsecase usecase.StateUsecase
	userUsecase  usecase.UserUsecase
	jwtSvc       *jwt.Service
	validate     *validator.Validate
}

// NewHandler создает новый экземпляр HTTP обработчика
func NewHandler(
	bookUsecase usecase.BookUsecase,
	tagUsecase usecase.TagUsecase,
	stateUsecase usecase.StateUsecase,
	userUsecase usecase.UserUsecase,
	jwtSvc *jwt.Service,
) *Handler {
	return &Handler{
		bookUsecase:  bookUsecase,
		tagUsecase:   tagUsecase,
		stateUsecase: stateUsecase,
		userUsecase:  userUsecase,
		jwtSvc:       jwtSvc,
		validate:     validator.New(),
	}
}

// InitRoutes инициализирует маршруты API
func (h *Handler) InitRoutes(r chi.Router) {
	// Swagger документация
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // Путь к swagger.json
	))

	// API маршруты
	r.Route("/api/v1", func(r chi.Router) {
		// Группа маршрутов для тегов
		r.Route("/tags", func(r chi.Router) {
			r.Post("/", h.createTag)
			r.Get("/", h.getAllTags)
			r.Get("/{id}", h.getTagByID)
			r.Get("/popular", h.getPopularTags)
			r.Delete("/{id}", h.deleteTag)
			r.Put("/{id}", h.updateTag)
		})

		// Группа маршрутов для книг
		r.Route("/books", func(r chi.Router) {
			r.Post("/", h.createBook)
			r.Get("/{id}", h.getBookByID)
			r.Get("/search", h.searchBooksByTags)
			r.Post("/{id}/tags", h.addTagsToBook)
			r.Put("/{id}", h.updateBook)
			r.Patch("/{id}/state", h.updateBookState)
			r.Delete("/{id}", h.deleteBook)
			r.Get("/", h.getAllBooks)
		})

		// Группа маршрутов для состояний
		r.Route("/states", func(r chi.Router) {
			r.Post("/", h.createState)
			r.Get("/", h.getAllStates)
			r.Get("/{id}", h.getStateByID)
			r.Put("/{id}", h.updateState)
			r.Delete("/{id}", h.deleteState)
		})
	})
}

// @Summary Create tag
// @Tags Tags
// @Accept json
// @Produce json
// @Param tag body tag.CreateTagDTO true "Tag details"
// @Success 201 {object} tag.Tag
// @Failure 400,401,500 {object} ErrorResponse
// @Security Bearer
// @Router /api/v1/tags [post]
// @x-tryItOutEnabled true
// @x-validateRequest true
// @x-validateResponse true
func (h *Handler) createTag(w http.ResponseWriter, r *http.Request) {
	var dto tag.CreateTagDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		logger.Error("Failed to decode request body", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the DTO
	if err := h.validate.Struct(dto); err != nil {
		logger.Error("Validation failed", err)
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Create new tag
	newTag := &tag.Tag{
		Name: dto.Name,
	}

	// Save tag
	if err := h.tagUsecase.CreateTag(newTag); err != nil {
		logger.Error("Failed to create tag", err)
		http.Error(w, "Failed to create tag: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTag)
}

// @Summary Get tag
// @Tags Tags
// @Produce json
// @Param id path int true "Tag ID"
// @Success 200 {object} tag.Tag
// @Failure 400,404 {object} ErrorResponse
// @Router /api/v1/tags/{id} [get]
func (h *Handler) getTagByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		logger.Error("Invalid tag ID", err)
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
		return
	}

	tag, err := h.tagUsecase.GetTagByID(uint(id))
	if err != nil {
		logger.Error("Failed to get tag", err)
		http.Error(w, "Tag not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tag)
}

// @Summary Get popular tags
// @Tags Tags
// @Produce json
// @Param limit query int false "Number of tags"
// @Success 200 {array} tag.Tag
// @Router /api/v1/tags/popular [get]
func (h *Handler) getPopularTags(w http.ResponseWriter, r *http.Request) {
	limit := 10 // Значение по умолчанию
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	tags, err := h.tagUsecase.GetPopularTags(limit)
	if err != nil {
		logger.Error("Failed to get popular tags", err)
		http.Error(w, "Failed to get popular tags", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}

// @Summary Create a new book
// @Description Create a new book with the given details
// @Tags books
// @Accept json
// @Produce json
// @Param book body book.CreateBookDTO true "Book details"
// @Success 201 {object} book.Book
// @Failure 400,401,500 {object} ErrorResponse
// @Security Bearer
// @Router /api/v1/books [post]
// @x-tryItOutEnabled true
// @x-validateRequest true
// @x-validateResponse true
func (h *Handler) createBook(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userID, ok := r.Context().Value(UserIDKey).(uint)
	if !ok {
		logger.Error("Failed to get user ID from context", errors.New("user ID not found in context"))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var dto book.CreateBookDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		logger.Error("Failed to decode request body", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(dto); err != nil {
		logger.Error("Validation failed", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем существование состояния
	state, err := h.stateUsecase.GetByID(uint(dto.StateID))
	if err != nil {
		logger.Error("Invalid state ID", err)
		http.Error(w, "Invalid state ID", http.StatusBadRequest)
		return
	}

	// Проверяем существование всех тегов
	for _, tagID := range dto.TagIDs {
		_, err := h.tagUsecase.GetTagByID(uint(tagID))
		if err != nil {
			logger.Error("Invalid tag ID", err)
			http.Error(w, fmt.Sprintf("Invalid tag ID: %d", tagID), http.StatusBadRequest)
			return
		}
	}

	// Create new book
	newBook := &book.Book{
		Title:       dto.Title,
		Author:      dto.Author,
		Description: dto.Description,
		Photos:      dto.Photos,
		UserID:      userID,
		StateID:     state.ID,
	}

	// Convert tag IDs to uint
	tagIDs := make([]uint, len(dto.TagIDs))
	for i, id := range dto.TagIDs {
		tagIDs[i] = uint(id)
	}

	// Save book
	if err := h.bookUsecase.CreateBook(newBook, tagIDs); err != nil {
		logger.Error("Failed to create book", err)
		http.Error(w, "Failed to create book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newBook)
}

// @Summary Update a book
// @Description Update book details
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param book body book.UpdateBookDTO true "Book details"
// @Success 200 {object} book.Book
// @Failure 400,401,404,500 {object} ErrorResponse
// @Security Bearer
// @Router /api/v1/books/{id} [put]
// @x-tryItOutEnabled true
// @x-validateRequest true
// @x-validateResponse true
func (h *Handler) updateBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		logger.Error("Invalid book ID", err)
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var dto book.UpdateBookDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		logger.Error("Failed to decode request body", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the DTO
	if err := h.validate.Struct(dto); err != nil {
		logger.Error("Validation failed", err)
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Update book
	updatedBook, err := h.bookUsecase.UpdateBook(uint(id), &dto)
	if err != nil {
		logger.Error("Failed to update book", err)
		http.Error(w, "Failed to update book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedBook)
}

// @Summary Get book by ID
// @Description Get book information by ID
// @Tags books
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} book.Book
// @Router /api/v1/books/{id} [get]
func (h *Handler) getBookByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		logger.Error("Invalid book ID", err)
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	book, err := h.bookUsecase.GetBookByID(uint(id))
	if err != nil {
		logger.Error("Failed to get book", err)
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// @Summary Search books by tags
// @Description Search books by tag IDs
// @Tags books
// @Produce json
// @Param tagIds query []int true "Tag IDs"
// @Success 200 {array} book.Book
// @Router /api/v1/books/search [get]
func (h *Handler) searchBooksByTags(w http.ResponseWriter, r *http.Request) {
	tagIDsStr := r.URL.Query()["tagIds"]
	if len(tagIDsStr) == 0 {
		http.Error(w, "No tag IDs provided", http.StatusBadRequest)
		return
	}

	tagIDs := make([]uint, len(tagIDsStr))
	for i, idStr := range tagIDsStr {
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			logger.Error("Invalid tag ID", err)
			http.Error(w, "Invalid tag ID", http.StatusBadRequest)
			return
		}
		tagIDs[i] = uint(id)
	}

	books, err := h.bookUsecase.GetBooksByTags(tagIDs)
	if err != nil {
		logger.Error("Failed to search books", err)
		http.Error(w, "Failed to search books", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// @Summary Add tags to book
// @Description Add tags to an existing book
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param tagIds body []int true "Tag IDs"
// @Success 200 {object} book.Book
// @Failure 400,401,404,500 {object} ErrorResponse
// @Security Bearer
// @Router /api/v1/books/{id}/tags [post]
func (h *Handler) addTagsToBook(w http.ResponseWriter, r *http.Request) {
	bookID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		logger.Error("Invalid book ID", err)
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var tagIDsInt []int64
	if err := json.NewDecoder(r.Body).Decode(&tagIDsInt); err != nil {
		logger.Error("Failed to decode request body", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tagIDs := make([]uint, len(tagIDsInt))
	for i, id := range tagIDsInt {
		tagIDs[i] = uint(id)
	}

	if err := h.bookUsecase.AddTagsToBook(uint(bookID), tagIDs); err != nil {
		logger.Error("Failed to add tags to book", err)
		http.Error(w, "Failed to add tags to book", http.StatusInternalServerError)
		return
	}

	book, err := h.bookUsecase.GetBookByID(uint(bookID))
	if err != nil {
		logger.Error("Failed to get updated book", err)
		http.Error(w, "Failed to get updated book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// @Summary Create a new state
// @Description Create a new book state
// @Tags states
// @Accept json
// @Produce json
// @Param state body state.CreateStateDTO true "State object"
// @Success 201 {object} state.State
// @Failure 400,401,500 {object} ErrorResponse
// @Security Bearer
// @Router /api/v1/states [post]
func (h *Handler) createState(w http.ResponseWriter, r *http.Request) {
	var s state.State
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		logger.Error("Failed to decode request body", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.stateUsecase.Create(&s); err != nil {
		logger.Error("Failed to create state", err)
		http.Error(w, "Failed to create state", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

// @Summary Get all states
// @Description Get list of all book states
// @Tags states
// @Produce json
// @Success 200 {array} state.State
// @Router /api/v1/states [get]
func (h *Handler) getAllStates(w http.ResponseWriter, r *http.Request) {
	states, err := h.stateUsecase.GetAll()
	if err != nil {
		logger.Error("Failed to get states", err)
		http.Error(w, "Failed to get states", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(states)
}

// @Summary Get state by ID
// @Description Get a book state by ID
// @Tags states
// @Produce json
// @Param id path int true "State ID"
// @Success 200 {object} state.State
// @Router /api/v1/states/{id} [get]
func (h *Handler) getStateByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		logger.Error("Invalid state ID", err)
		http.Error(w, "Invalid state ID", http.StatusBadRequest)
		return
	}

	state, err := h.stateUsecase.GetByID(uint(id))
	if err != nil {
		logger.Error("Failed to get state", err)
		http.Error(w, "State not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

// @Summary Update state
// @Description Update a book state
// @Tags states
// @Accept json
// @Produce json
// @Param id path int true "State ID"
// @Param state body state.UpdateStateDTO true "State object"
// @Success 200 {object} state.State
// @Failure 400,401,404,500 {object} ErrorResponse
// @Security Bearer
// @Router /api/v1/states/{id} [put]
func (h *Handler) updateState(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		logger.Error("Invalid state ID", err)
		http.Error(w, "Invalid state ID", http.StatusBadRequest)
		return
	}

	var s state.State
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		logger.Error("Failed to decode request body", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	s.ID = uint(id)
	if err := h.stateUsecase.Update(&s); err != nil {
		logger.Error("Failed to update state", err)
		http.Error(w, "Failed to update state", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

// @Summary Delete state
// @Description Delete a book state
// @Tags states
// @Param id path int true "State ID"
// @Success 204 "No Content"
// @Failure 400,401,404,500 {object} ErrorResponse
// @Security Bearer
// @Router /api/v1/states/{id} [delete]
func (h *Handler) deleteState(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		logger.Error("Invalid state ID", err)
		http.Error(w, "Invalid state ID", http.StatusBadRequest)
		return
	}

	if err := h.stateUsecase.Delete(uint(id)); err != nil {
		logger.Error("Failed to delete state", err)
		http.Error(w, "Failed to delete state", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Delete book
// @Description Delete a book by ID
// @Tags books
// @Param id path int true "Book ID"
// @Success 204 "No Content"
// @Failure 400,401,404,500 {object} ErrorResponse
// @Security Bearer
// @Router /api/v1/books/{id} [delete]
func (h *Handler) deleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	if err := h.bookUsecase.DeleteBook(uint(id)); err != nil {
		http.Error(w, "Failed to delete book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Delete tag
// @Tags Tags
// @Param id path int true "Tag ID"
// @Success 204 "No Content"
// @Failure 400,401,404,500 {object} ErrorResponse
// @Security Bearer
// @Router /api/v1/tags/{id} [delete]
func (h *Handler) deleteTag(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
		return
	}

	if err := h.tagUsecase.DeleteTag(uint(id)); err != nil {
		http.Error(w, "Failed to delete tag", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Get all tags
// @Tags Tags
// @Produce json
// @Success 200 {array} tag.Tag
// @Router /api/v1/tags [get]
func (h *Handler) getAllTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.tagUsecase.GetAllTags()
	if err != nil {
		logger.Error("Failed to get tags", err)
		http.Error(w, "Failed to get tags", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}

// @Summary Update book state
// @Description Update the state of a book
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param state body book.UpdateBookStateDTO true "New state"
// @Success 200 {object} book.Book
// @Failure 400,401,404,500 {object} ErrorResponse
// @Security Bearer
// @Router /api/v1/books/{id}/state [patch]
func (h *Handler) updateBookState(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		logger.Error("Invalid book ID", err)
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var dto book.UpdateBookStateDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		logger.Error("Failed to decode request body", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(dto); err != nil {
		logger.Error("Validation failed", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book, err := h.bookUsecase.UpdateBookState(uint(id), uint(dto.StateID))
	if err != nil {
		logger.Error("Failed to update book state", err)
		http.Error(w, "Failed to update book state", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// @Summary Get all books
// @Description Get a list of all books with pagination
// @Tags books
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Items per page (default: 10, max: 100)"
// @Success 200 {object} map[string]interface{} "Returns books and pagination info"
// @Router /api/v1/books [get]
func (h *Handler) getAllBooks(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры пагинации
	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 10
	if pageSizeStr := r.URL.Query().Get("pageSize"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	// Получаем книги с пагинацией
	books, total, err := h.bookUsecase.GetAllBooks(page, pageSize)
	if err != nil {
		logger.Error("Failed to get books", err)
		http.Error(w, "Failed to get books", http.StatusInternalServerError)
		return
	}

	// Формируем ответ с информацией о пагинации
	response := map[string]interface{}{
		"books": books,
		"pagination": map[string]interface{}{
			"total":      total,
			"page":       page,
			"pageSize":   pageSize,
			"totalPages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary Get user books
// @Description Get paginated list of user's books
// @Tags books
// @Produce json
// @Param id path int true "User ID"
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Items per page (default: 10, max: 100)"
// @Success 200 {object} map[string]interface{} "Returns books and pagination info"
// @Security Bearer
// @Router /api/v1/users/{id}/books [get]
func (h *Handler) getUserBooks(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		logger.Error("Invalid user ID", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	page := 1
	pageSize := 10

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if sizeStr := r.URL.Query().Get("size"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 && s <= 100 {
			pageSize = s
		}
	}

	books, total, err := h.bookUsecase.GetUserBooks(uint(userID), page, pageSize)
	if err != nil {
		logger.Error("Failed to get user books", err)
		http.Error(w, "Failed to get user books", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"books": books,
		"total": total,
	})
}

// @Summary Update tag
// @Tags Tags
// @Accept json
// @Produce json
// @Param id path int true "Tag ID"
// @Param tag body tag.UpdateTagDTO true "Tag details"
// @Success 200 {object} tag.Tag
// @Failure 400,401,404,500 {object} ErrorResponse
// @Security Bearer
// @Router /api/v1/tags/{id} [put]
func (h *Handler) updateTag(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		logger.Error("Invalid tag ID", err)
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
		return
	}

	var dto tag.UpdateTagDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		logger.Error("Failed to decode request body", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(dto); err != nil {
		logger.Error("Validation failed", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedTag, err := h.tagUsecase.UpdateTag(uint(id), &dto)
	if err != nil {
		logger.Error("Failed to update tag", err)
		http.Error(w, "Failed to update tag", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTag)
}

// @Summary Refresh access token
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} jwt.TokenPair
// @Failure 401,500 {object} ErrorResponse
// @Security RefreshToken
// @Router /api/v1/auth/refresh [post]
func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.Header.Get("X-Refresh-Token")
	if refreshToken == "" {
		http.Error(w, "Refresh token is required", http.StatusUnauthorized)
		return
	}

	// Получаем пользователя по refresh token
	user, err := h.jwtSvc.ValidateRefreshToken(refreshToken)
	if err != nil {
		logger.Error("Failed to validate refresh token", err)
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Генерируем новую пару токенов
	tokenPair, err := h.jwtSvc.GenerateTokenPair(user)
	if err != nil {
		logger.Error("Failed to generate token pair", err)
		http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenPair)
}
