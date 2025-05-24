package http

import (
	"booktrading/internal/domain/user"
	"booktrading/internal/pkg/logger"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// @Summary Register new user
// @Description Register a new user in the system
// @Tags users
// @Accept json
// @Produce json
// @Param user body user.CreateUserDTO true "User registration data"
// @Success 201 {object} user.User
// @Router /api/v1/users/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var dto user.CreateUserDTO
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

	user, err := h.userUsecase.Register(&dto)
	if err != nil {
		logger.Error("Failed to register user", err)
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// @Summary Login user
// @Description Login user and get JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body user.LoginDTO true "Login credentials"
// @Success 200 {object} user.TokenResponse
// @Router /api/v1/users/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var dto user.LoginDTO
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

	tokenPair, err := h.userUsecase.Login(&dto)
	if err != nil {
		logger.Error("Failed to login", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenPair)
}

// @Summary Get user by ID
// @Description Get user information by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} user.User
// @Router /api/v1/users/{id} [get]
func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		logger.Error("Invalid user ID", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userUsecase.GetByID(uint(id))
	if err != nil {
		logger.Error("Failed to get user", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// @Summary Update user
// @Description Update user information
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body user.UpdateUserDTO true "User update data"
// @Success 200 {object} user.User
// @Security BearerAuth
// @Router /api/v1/users/{id} [put]
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		logger.Error("Invalid user ID", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var dto user.UpdateUserDTO
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

	updatedUser, err := h.userUsecase.Update(uint(id), &dto)
	if err != nil {
		logger.Error("Failed to update user", err)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

// @Summary Delete user
// @Description Delete user by ID
// @Tags users
// @Param id path int true "User ID"
// @Success 204 "No Content"
// @Security BearerAuth
// @Router /api/v1/users/{id} [delete]
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		logger.Error("Invalid user ID", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.userUsecase.Delete(uint(id)); err != nil {
		logger.Error("Failed to delete user", err)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Get all users
// @Description Get a paginated list of users
// @Tags users
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Success 200 {object} map[string]interface{} "Users and total count"
// @Router /api/v1/users [get]
func (h *Handler) getAllUsers(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры пагинации из query string
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	page := 1
	pageSize := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	users, total, err := h.userUsecase.GetAll(page, pageSize)
	if err != nil {
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"users": users,
		"total": total,
		"page":  page,
		"size":  pageSize,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
