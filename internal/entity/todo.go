package entity

import "time"

type TodoStatus string

const (
	TodoPending   TodoStatus = "pending"
	TodoCompleted TodoStatus = "completed"
)

type Todo struct {
	ID          int        `json:"id"`
	UserID      int        `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TodoStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
