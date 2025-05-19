package http

import (
	"booktrading/internal/domain/user"
	"booktrading/internal/usecase"
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
func (h *Handler) registerUser(w http.ResponseWriter, r *http.Request) {
	var dto user.CreateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(dto); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	newUser, err := h.userUsecase.Register(&dto)
	if err != nil {
		if err == usecase.ErrUserAlreadyExists {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

// @Summary Login user
// @Description Login user and get JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body user.LoginDTO true "Login credentials"
// @Success 200 {object} user.TokenResponse
// @Router /api/v1/users/login [post]
func (h *Handler) loginUser(w http.ResponseWriter, r *http.Request) {
	var dto user.LoginDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(dto); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.userUsecase.Login(&dto)
	if err != nil {
		if err == usecase.ErrInvalidCredentials {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Failed to login", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

// @Summary Get user by ID
// @Description Get user information by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} user.User
// @Router /api/v1/users/{id} [get]
func (h *Handler) getUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userUsecase.GetByID(uint(id))
	if err != nil {
		if err == usecase.ErrUserNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
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
// @Router /api/v1/users/{id} [put]
func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var dto user.UpdateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(dto); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	updatedUser, err := h.userUsecase.Update(uint(id), &dto)
	if err != nil {
		if err == usecase.ErrUserNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
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
// @Router /api/v1/users/{id} [delete]
func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.userUsecase.Delete(uint(id)); err != nil {
		if err == usecase.ErrUserNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
