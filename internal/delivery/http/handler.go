package http

import (
	"booktrading/internal/domain/book"
	"booktrading/internal/domain/response"
	"booktrading/internal/domain/state"
	"booktrading/internal/domain/tag"
	"booktrading/internal/domain/user"
	"booktrading/internal/pkg/logger"
	"booktrading/internal/pkg/validator"
	"booktrading/internal/usecase"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
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
	bookUsecase  usecase.BookUseCase
	tagUsecase   usecase.TagUseCase
	stateUsecase usecase.StateUseCase
	userUsecase  usecase.UserUseCase
	validate     *validator.Validate
}

// error отправляет ответ с ошибкой
// @Description Отправляет ответ с ошибкой
// @Param code body int true "HTTP status code"
// @Param message body string true "Error message"
// @Success 200 {object} ErrorResponse
func (h *Handler) error(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

// respond отправляет успешный ответ
// @Description Отправляет успешный ответ
// @Param code body int true "HTTP status code"
// @Param data body object true "Response data"
// @Success 200 {object} object
func (h *Handler) respond(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

// NewHandler создает новый экземпляр HTTP обработчика
func NewHandler(
	bookUsecase usecase.BookUseCase,
	tagUsecase usecase.TagUseCase,
	stateUsecase usecase.StateUseCase,
	userUsecase usecase.UserUseCase,
) *Handler {
	return &Handler{
		bookUsecase:  bookUsecase,
		tagUsecase:   tagUsecase,
		stateUsecase: stateUsecase,
		userUsecase:  userUsecase,
		validate:     validator.New(),
	}
}

// @Summary Create new tag
// @Description Create a new tag with the provided name and optional photo
// @Tags Tags
// @Accept json
// @Produce json
// @Param tag body tag.CreateTagDTO true "Tag data"
// @Success 201 {object} tag.Tag
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security Bearer
// @Router /api/v1/tags [post]
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
		Name:  dto.Name,
		Photo: dto.Photo,
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

// @Summary Get tag by ID
// @Description Get tag information by its ID
// @Tags Tags
// @Produce json
// @Param id path int true "Tag ID"
// @Success 200 {object} tag.Tag
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
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
// @Description Get list of popular tags with optional limit
// @Tags Tags
// @Produce json
// @Param limit query int false "Number of tags to return (default: 10)"
// @Success 200 {array} tag.Tag
// @Failure 500 {object} ErrorResponse
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

// @Summary Create new book
// @Description Create a new book with the provided details
// @Tags Books
// @Accept json
// @Produce json
// @Param book body book.CreateBookDTO true "Book data"
// @Success 201 {object} book.Book "Created book"
// @Failure 400 {object} ErrorResponse "Invalid request data"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Security Bearer
// @Router /api/v1/books [post]
// @Example {json} Request:
//
//	{
//	  "title": "Война и мир",
//	  "author": "Лев Толстой",
//	  "description": "Роман-эпопея, описывающий русское общество в эпоху войн против Наполеона",
//	  "photos": [
//	    {
//	      "photo_url": "data:image/jpeg;base64,/9j/4AAQSkZJRg...",
//	      "is_main": true
//	    }
//	  ],
//	  "user_id": 1,
//	  "state_id": 1,
//	  "tag_ids": [1, 2, 3]
//	}
//
// @Example {json} Success Response:
//
//	{
//	  "id": 1,
//	  "title": "Война и мир",
//	  "author": "Лев Толстой",
//	  "description": "Роман-эпопея, описывающий русское общество в эпоху войн против Наполеона",
//	  "photos": [
//	    {
//	      "id": 1,
//	      "book_id": 1,
//	      "photo_url": "data:image/jpeg;base64,/9j/4AAQSkZJRg...",
//	      "is_main": true,
//	      "created_at": "2024-03-20T12:00:00Z",
//	      "updated_at": "2024-03-20T12:00:00Z"
//	    }
//	  ],
//	  "user_id": 1,
//	  "state_id": 1,
//	  "created_at": "2024-03-20T12:00:00Z",
//	  "updated_at": "2024-03-20T12:00:00Z"
//	}
func (h *Handler) createBook(w http.ResponseWriter, r *http.Request) {
	// Get claims from context
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		logger.Error("Failed to get claims from context", err)
		http.Error(w, "Authentication failed: Token not found", http.StatusUnauthorized)
		return
	}

	// Get user ID from claims
	userID, ok := claims["user_id"].(float64)
	if !ok {
		logger.Error("User ID not found in claims", nil)
		http.Error(w, "Authentication failed: User ID not found in token claims", http.StatusUnauthorized)
		return
	}

	var dto book.CreateBookDTO
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

	// Create new book with user ID from token
	newBook := dto.ToBook()
	newBook.UserID = uint(userID)

	// Save book and associate tags
	if err := h.bookUsecase.CreateBook(newBook, dto.TagIDs); err != nil {
		logger.Error("Failed to create book", err)
		http.Error(w, "Failed to create book: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create photos if they exist
	if len(dto.Photos) > 0 {
		for _, photoData := range dto.Photos {
			photo := &book.BookPhoto{
				BookID:   newBook.ID,
				PhotoURL: photoData.PhotoURL,
				IsMain:   photoData.IsMain,
			}
			if err := h.bookUsecase.CreatePhoto(photo); err != nil {
				logger.Error("Failed to create photo", err)
				http.Error(w, "Failed to create photo: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newBook)
}

// @Summary Get book by ID
// @Description Get book information by its ID
// @Tags Books
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} book.Book
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
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

// @Summary Update book
// @Description Update existing book information. Only the book owner can update it.
// @Tags Books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param book body book.UpdateBookDTO true "Updated book data"
// @Success 200 {object} book.Book "Updated book information"
// @Failure 400 {object} ErrorResponse "Invalid request data or validation failed"
// @Failure 401 {object} ErrorResponse "Unauthorized - Invalid or missing token"
// @Failure 403 {object} ErrorResponse "Forbidden - User is not the book owner"
// @Failure 404 {object} ErrorResponse "Book not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Security Bearer
// @Router /api/v1/books/{id} [put]
// @Example {json} Request:
//
//	{
//	  "title": "Updated Book Title",
//	  "author": "Updated Author",
//	  "description": "Updated book description",
//	  "photos": [
//	    {
//	      "photo_url": "data:image/jpeg;base64,...",
//	      "is_main": true
//	    }
//	  ],
//	  "tag_ids": [1, 2, 3]
//	}
func (h *Handler) updateBook(w http.ResponseWriter, r *http.Request) {
	// Get book ID from URL
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		logger.Error("Invalid book ID", err)
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	// Get claims from context
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		logger.Error("Failed to get claims from context", err)
		http.Error(w, "Authentication failed: Token not found", http.StatusUnauthorized)
		return
	}

	// Get user ID from claims
	userID, ok := claims["user_id"].(float64)
	if !ok {
		logger.Error("User ID not found in claims", nil)
		http.Error(w, "Authentication failed: User ID not found in token claims", http.StatusUnauthorized)
		return
	}

	// Get existing book
	existingBook, err := h.bookUsecase.GetBookByID(uint(id))
	if err != nil {
		logger.Error("Failed to get book", err)
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	// Check if user owns the book
	if existingBook.UserID != uint(userID) {
		logger.Error("User does not own the book", nil)
		http.Error(w, "You don't have permission to update this book", http.StatusForbidden)
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

	// Update book fields
	existingBook.UpdateFromDTO(&dto)

	// Update book
	if err := h.bookUsecase.UpdateBook(existingBook, dto.TagIDs); err != nil {
		logger.Error("Failed to update book", err)
		http.Error(w, "Failed to update book: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Update photos if they exist
	if len(dto.Photos) > 0 {
		// Delete existing photos
		if err := h.bookUsecase.DeletePhotos(existingBook.ID); err != nil {
			logger.Error("Failed to delete existing photos", err)
			http.Error(w, "Failed to delete existing photos: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Create new photos
		for i, photoURL := range dto.Photos {
			photo := &book.BookPhoto{
				BookID:   existingBook.ID,
				PhotoURL: photoURL,
				IsMain:   i == 0, // Первая фотография - главная
			}
			if err := h.bookUsecase.CreatePhoto(photo); err != nil {
				logger.Error("Failed to create photo", err)
				http.Error(w, "Failed to create photo: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingBook)
}

// @Summary Delete book
// @Description Delete book by ID
// @Tags Books
// @Param id path int true "Book ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
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

// @Summary Search books by tags
// @Description Search books by provided tag IDs
// @Tags Books
// @Produce json
// @Param tagIds query []int true "Tag IDs"
// @Success 200 {array} book.Book
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
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
// @Description Add new tags to an existing book
// @Tags Books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param tagIds body []int true "Tag IDs to add"
// @Success 200 {object} book.Book
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
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

// @Summary Create new state
// @Description Create a new book state
// @Tags States
// @Accept json
// @Produce json
// @Param state body state.CreateStateDTO true "State data"
// @Success 201 {object} state.State
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
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
// @Tags States
// @Produce json
// @Success 200 {array} state.State
// @Failure 500 {object} ErrorResponse
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
// @Description Get book state information by its ID
// @Tags States
// @Produce json
// @Param id path int true "State ID"
// @Success 200 {object} state.State
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
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
// @Description Update existing book state information
// @Tags States
// @Accept json
// @Produce json
// @Param id path int true "State ID"
// @Param state body state.UpdateStateDTO true "Updated state data"
// @Success 200 {object} state.State
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
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
// @Description Delete book state by ID
// @Tags States
// @Param id path int true "State ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
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

// @Summary Get all books
// @Description Get paginated list of all books
// @Tags Books
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Items per page (default: 10, max: 100)"
// @Success 200 {object} map[string]interface{} "Returns books and pagination info"
// @Failure 500 {object} ErrorResponse
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
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Items per page (default: 10, max: 100)"
// @Success 200 {object} map[string]interface{} "Returns books and pagination info"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
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
// @Description Update existing tag information including optional photo
// @Tags Tags
// @Accept json
// @Produce json
// @Param id path int true "Tag ID"
// @Param tag body tag.UpdateTagDTO true "Updated tag data"
// @Success 200 {object} tag.Tag
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
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
// @Success 200 {object} response.TokenResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security RefreshToken
// @Router /api/v1/auth/refresh [post]
func (h *Handler) refreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.Header.Get("X-Refresh-Token")
	if refreshToken == "" {
		http.Error(w, "Refresh token is required", http.StatusUnauthorized)
		return
	}

	// Получаем пользователя по refresh token
	user, err := h.userUsecase.ValidateRefreshToken(refreshToken)
	if err != nil {
		logger.Error("Failed to validate refresh token", err)
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Генерируем новую пару токенов
	tokenPair, err := h.userUsecase.GenerateTokenPair(user)
	if err != nil {
		logger.Error("Failed to generate token pair", err)
		http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenPair)
}

// @Summary Get all users
// @Description Get paginated list of all users
// @Tags Users
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Items per page (default: 10, max: 100)"
// @Success 200 {object} map[string]interface{} "Returns users and pagination info"
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security Bearer
// @Router /api/v1/users [get]
func (h *Handler) getAllUsers(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры пагинации из запроса
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

	// Call the usecase to get all users
	users, total, err := h.userUsecase.GetAll(page, pageSize)
	if err != nil {
		logger.Error("Failed to get all users", err)
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	// Format the response with pagination info
	response := map[string]interface{}{
		"users": users,
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

// @Summary Get user by ID
// @Description Get user information by ID
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} user.User
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security Bearer
// @Router /api/v1/users/{id} [get]
func (h *Handler) getUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	u, err := h.userUsecase.GetByID(uint(id))
	if err != nil {
		logger.Error("Failed to get user by ID", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

// @Summary Update user
// @Description Update existing user information
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body user.UpdateUserDTO true "Updated user data"
// @Success 200 {object} user.User
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security Bearer
// @Router /api/v1/users/{id} [put]
func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из URL
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req user.UpdateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedUser, err := h.userUsecase.Update(uint(id), &req)
	if err != nil {
		logger.Error("Failed to update user", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

// @Summary Delete user
// @Description Delete user by ID
// @Tags Users
// @Param id path int true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security Bearer
// @Router /api/v1/users/{id} [delete]
func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from URL parameter
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Error("Invalid user ID for deletion", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Call the usecase to delete the user
	if err := h.userUsecase.Delete(uint(id)); err != nil {
		logger.Error("Failed to delete user", err)
		if errors.Is(err, usecase.ErrUserNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Update book state
// @Description Update the state of an existing book. Only the book owner can update its state.
// @Tags Books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param state body book.UpdateBookStateDTO true "New state data"
// @Success 200 {object} book.Book "Updated book with new state"
// @Failure 400 {object} ErrorResponse "Invalid request data or validation failed"
// @Failure 401 {object} ErrorResponse "Unauthorized - Invalid or missing token"
// @Failure 403 {object} ErrorResponse "Forbidden - User is not the book owner"
// @Failure 404 {object} ErrorResponse "Book not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Security Bearer
// @Router /api/v1/books/{id}/state [patch]
// @Example {json} Request:
//
//	{
//	  "state_id": 2
//	}
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

// @Summary Register new user
// @Description Register a new user with the provided credentials
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body user.CreateUserDTO true "User registration data"
// @Success 201 {object} user.User
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/auth/register [post]
func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var dto user.CreateUserDTO
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

	// Register new user
	newUser, err := h.userUsecase.Register(&dto)
	if err != nil {
		logger.Error("Failed to register user", err)
		if errors.Is(err, usecase.ErrUserAlreadyExists) {
			http.Error(w, "User already exists", http.StatusConflict)
		} else {
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

// @Summary Login user
// @Description Authenticate user and return JWT tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body user.LoginDTO true "Данные для входа"
// @Success 200 {object} response.LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var dto user.LoginDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tokenResponse, userID, err := h.userUsecase.Login(&dto)
	if err != nil {
		if err == usecase.ErrInvalidCredentials {
			h.error(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		h.error(w, http.StatusInternalServerError, "failed to login")
		return
	}

	response := response.LoginResponse{
		Token:        tokenResponse.Token,
		RefreshToken: tokenResponse.RefreshToken,
		UserID:       userID,
	}

	h.respond(w, http.StatusOK, response)
}

// @Summary Get all tags
// @Description Get list of all tags
// @Tags Tags
// @Produce json
// @Success 200 {array} tag.Tag
// @Failure 500 {object} ErrorResponse
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

// @Summary Delete tag
// @Description Delete tag by ID
// @Tags Tags
// @Param id path int true "Tag ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
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

// @Summary Logout user
// @Description Logout user and invalidate refresh token. Requires a valid refresh token in the X-Refresh-Token header.
// @Tags Auth
// @Accept json
// @Produce json
// @Header 204 {string} X-Refresh-Token "Refresh token to invalidate"
// @Success 204 "No Content - Logout successful"
// @Failure 400 {object} ErrorResponse "Refresh token is required"
// @Failure 401 {object} ErrorResponse "Invalid refresh token"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/auth/logout [post]
// @Example {json} Request Header:
//
//	X-Refresh-Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.Header.Get("X-Refresh-Token")
	if refreshToken == "" {
		h.error(w, http.StatusBadRequest, "Refresh token is required")
		return
	}

	if err := h.userUsecase.Logout(refreshToken); err != nil {
		logger.Error("Failed to logout user", err)
		h.error(w, http.StatusInternalServerError, "Failed to logout")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
