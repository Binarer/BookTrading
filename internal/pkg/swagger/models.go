package swagger

import "time"

// DeletedAt представляет время удаления записи
// @Description Время удаления записи для мягкого удаления
type DeletedAt struct {
	Time  time.Time
	Valid bool
}

// Base содержит общие поля для всех таблиц
// @Description Базовая структура для всех моделей
type Base struct {
	// @Description Уникальный идентификатор
	// @example 1
	ID        uint      `json:"id"`
	
	// @Description Время создания записи
	// @example 2025-04-28T12:00:00Z
	CreatedAt time.Time `json:"created_at"`
	
	// @Description Время последнего обновления записи
	// @example 2025-04-28T12:00:00Z
	UpdatedAt time.Time `json:"updated_at"`
	
	// @Description Время удаления записи (для мягкого удаления)
	// @example 2025-04-28T12:00:00Z
	DeletedAt DeletedAt `json:"deleted_at,omitempty"`
}

// Tag представляет тег книги
// @Description Модель тега для категоризации книг
type Tag struct {
	Base
	// @Description Название тега
	// @example fiction
	Name string `json:"name"`
}

// Book представляет сущность книги
// @Description Модель книги для системы обмена книгами
type Book struct {
	Base
	// @Description Название книги
	// @example The Great Gatsby
	Title       string `json:"title"`
	
	// @Description Автор книги
	// @example F. Scott Fitzgerald
	Author      string `json:"author"`
	
	// @Description Описание книги
	// @example A story of the fabulously wealthy Jay Gatsby and his love for the beautiful Daisy Buchanan.
	Description string `json:"description"`
	
	// @Description Цена книги
	// @example 19.99
	Price       float64 `json:"price"`
	
	// @Description ID состояния книги
	// @example 1
	StateID     uint `json:"state_id"`
	
	// @Description ID пользователя-владельца
	// @example 1
	UserID      uint `json:"user_id"`
}

// State представляет состояние книги
// @Description Модель состояния книги
type State struct {
	Base
	// @Description Название состояния
	// @example available
	Name string `json:"name"`
}

// User представляет пользователя системы
// @Description Модель пользователя системы обмена книгами
type User struct {
	Base
	// @Description Имя пользователя
	// @example john_doe
	Username     string    `json:"username"`
	
	// @Description Email пользователя
	// @example john@example.com
	Email        string    `json:"email"`
	
	// @Description Время создания аккаунта
	// @example 2025-04-28T12:00:00Z
	CreatedAt    time.Time `json:"created_at"`
	
	// @Description Время последнего обновления аккаунта
	// @example 2025-04-28T12:00:00Z
	UpdatedAt    time.Time `json:"updated_at"`
} 