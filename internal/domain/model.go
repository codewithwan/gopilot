package domain

import (
	"time"
)

type Todo struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	UserID      int64     `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateTodoRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=255"`
	Description string `json:"description" binding:"max=1000"`
}

type UpdateTodoRequest struct {
	Title       *string `json:"title" binding:"omitempty,min=1,max=255"`
	Description *string `json:"description" binding:"omitempty,max=1000"`
	Completed   *bool   `json:"completed"`
}

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=100"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
