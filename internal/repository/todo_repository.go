package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/islamyakin/otel-propagation-monorepo/internal/entity"
	"github.com/islamyakin/otel-propagation-monorepo/internal/model"
	"github.com/islamyakin/otel-propagation-monorepo/internal/model/converter"
)

type TodoRepository interface {
	Create(todo *entity.Todo) (*entity.Todo, error)
	GetByID(id int) (*entity.Todo, error)
	GetByUserID(userID int) ([]*entity.Todo, error)
	GetAll() ([]*entity.Todo, error)
	Update(todo *entity.Todo) (*entity.Todo, error)
	Delete(id int) error
	GetByIDAndUserID(id, userID int) (*entity.Todo, error)
}

type todoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) Create(todo *entity.Todo) (*entity.Todo, error) {
	query := `
		INSERT INTO todos (user_id, title, description, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, title, description, status, created_at, updated_at
	`

	now := time.Now()
	var todoModel model.TodoModel

	err := r.db.QueryRow(query, todo.UserID, todo.Title, todo.Description, string(todo.Status), now, now).
		Scan(&todoModel.ID, &todoModel.UserID, &todoModel.Title, &todoModel.Description, &todoModel.Status, &todoModel.CreatedAt, &todoModel.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}

	return converter.TodoModelToEntity(&todoModel), nil
}

func (r *todoRepository) GetByID(id int) (*entity.Todo, error) {
	query := `
		SELECT id, user_id, title, description, status, created_at, updated_at
		FROM todos
		WHERE id = $1
	`

	var todoModel model.TodoModel
	err := r.db.QueryRow(query, id).
		Scan(&todoModel.ID, &todoModel.UserID, &todoModel.Title, &todoModel.Description, &todoModel.Status, &todoModel.CreatedAt, &todoModel.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("todo not found")
		}
		return nil, fmt.Errorf("failed to get todo by id: %w", err)
	}

	return converter.TodoModelToEntity(&todoModel), nil
}

func (r *todoRepository) GetByUserID(userID int) ([]*entity.Todo, error) {
	query := `
		SELECT id, user_id, title, description, status, created_at, updated_at
		FROM todos
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get todos by user id: %w", err)
	}
	defer rows.Close()

	var todoModels []*model.TodoModel
	for rows.Next() {
		var todoModel model.TodoModel
		err := rows.Scan(&todoModel.ID, &todoModel.UserID, &todoModel.Title, &todoModel.Description, &todoModel.Status, &todoModel.CreatedAt, &todoModel.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan todo: %w", err)
		}
		todoModels = append(todoModels, &todoModel)
	}

	return converter.TodoModelsToEntities(todoModels), nil
}

func (r *todoRepository) GetAll() ([]*entity.Todo, error) {
	query := `
		SELECT id, user_id, title, description, status, created_at, updated_at
		FROM todos
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all todos: %w", err)
	}
	defer rows.Close()

	var todoModels []*model.TodoModel
	for rows.Next() {
		var todoModel model.TodoModel
		err := rows.Scan(&todoModel.ID, &todoModel.UserID, &todoModel.Title, &todoModel.Description, &todoModel.Status, &todoModel.CreatedAt, &todoModel.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan todo: %w", err)
		}
		todoModels = append(todoModels, &todoModel)
	}

	return converter.TodoModelsToEntities(todoModels), nil
}

func (r *todoRepository) Update(todo *entity.Todo) (*entity.Todo, error) {
	query := `
		UPDATE todos
		SET title = $2, description = $3, status = $4, updated_at = $5
		WHERE id = $1
		RETURNING id, user_id, title, description, status, created_at, updated_at
	`

	now := time.Now()
	var todoModel model.TodoModel

	err := r.db.QueryRow(query, todo.ID, todo.Title, todo.Description, string(todo.Status), now).
		Scan(&todoModel.ID, &todoModel.UserID, &todoModel.Title, &todoModel.Description, &todoModel.Status, &todoModel.CreatedAt, &todoModel.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("todo not found")
		}
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}

	return converter.TodoModelToEntity(&todoModel), nil
}

func (r *todoRepository) Delete(id int) error {
	query := `DELETE FROM todos WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("todo not found")
	}

	return nil
}

func (r *todoRepository) GetByIDAndUserID(id, userID int) (*entity.Todo, error) {
	query := `
		SELECT id, user_id, title, description, status, created_at, updated_at
		FROM todos
		WHERE id = $1 AND user_id = $2
	`

	var todoModel model.TodoModel
	err := r.db.QueryRow(query, id, userID).
		Scan(&todoModel.ID, &todoModel.UserID, &todoModel.Title, &todoModel.Description, &todoModel.Status, &todoModel.CreatedAt, &todoModel.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("todo not found")
		}
		return nil, fmt.Errorf("failed to get todo by id and user id: %w", err)
	}

	return converter.TodoModelToEntity(&todoModel), nil
}
