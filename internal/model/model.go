package model

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/islamyakin/otel-propagation-monorepo/internal/entity"
)

type UserModel struct {
	ID        int       `db:"id"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type TodoModel struct {
	ID          int       `db:"id"`
	UserID      int       `db:"user_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Status      string    `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// Register request/response models
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  entity.User `json:"user"`
}

type CreateTodoRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type UpdateTodoRequest struct {
	Title       *string            `json:"title"`
	Description *string            `json:"description"`
	Status      *entity.TodoStatus `json:"status"`
}

// Custom type for role that implements sql driver interfaces
type Role string

func (r Role) Value() (driver.Value, error) {
	return string(r), nil
}

func (r *Role) Scan(value interface{}) error {
	if value == nil {
		*r = ""
		return nil
	}
	switch s := value.(type) {
	case string:
		*r = Role(s)
	case []byte:
		*r = Role(s)
	default:
		return fmt.Errorf("cannot scan %T into Role", value)
	}
	return nil
}
