package usecase

import (
	"fmt"

	"github.com/islamyakin/otel-propagation-monorepo/internal/entity"
	"github.com/islamyakin/otel-propagation-monorepo/internal/model"
	"github.com/islamyakin/otel-propagation-monorepo/internal/repository"
)

type TodoUseCase interface {
	Create(userID int, req *model.CreateTodoRequest) (*entity.Todo, error)
	GetByUserID(userID int) ([]*entity.Todo, error)
	GetAll() ([]*entity.Todo, error) // Admin only
	GetByID(todoID, userID int, isAdmin bool) (*entity.Todo, error)
	Update(todoID, userID int, req *model.UpdateTodoRequest, isAdmin bool) (*entity.Todo, error)
	Delete(todoID, userID int, isAdmin bool) error
}

type todoUseCase struct {
	todoRepo repository.TodoRepository
}

func NewTodoUseCase(todoRepo repository.TodoRepository) TodoUseCase {
	return &todoUseCase{
		todoRepo: todoRepo,
	}
}

func (uc *todoUseCase) Create(userID int, req *model.CreateTodoRequest) (*entity.Todo, error) {
	todo := &entity.Todo{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Status:      entity.TodoPending,
	}

	createdTodo, err := uc.todoRepo.Create(todo)
	if err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}

	return createdTodo, nil
}

func (uc *todoUseCase) GetByUserID(userID int) ([]*entity.Todo, error) {
	todos, err := uc.todoRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get todos: %w", err)
	}

	return todos, nil
}

func (uc *todoUseCase) GetAll() ([]*entity.Todo, error) {
	todos, err := uc.todoRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all todos: %w", err)
	}

	return todos, nil
}

func (uc *todoUseCase) GetByID(todoID, userID int, isAdmin bool) (*entity.Todo, error) {
	var todo *entity.Todo
	var err error

	if isAdmin {
		// Admin can see any todo
		todo, err = uc.todoRepo.GetByID(todoID)
	} else {
		// User can only see their own todo
		todo, err = uc.todoRepo.GetByIDAndUserID(todoID, userID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}

	return todo, nil
}

func (uc *todoUseCase) Update(todoID, userID int, req *model.UpdateTodoRequest, isAdmin bool) (*entity.Todo, error) {
	var existingTodo *entity.Todo
	var err error

	if isAdmin {
		// Admin can update any todo
		existingTodo, err = uc.todoRepo.GetByID(todoID)
	} else {
		// User can only update their own todo
		existingTodo, err = uc.todoRepo.GetByIDAndUserID(todoID, userID)
	}

	if err != nil {
		return nil, fmt.Errorf("todo not found or access denied: %w", err)
	}

	// Update fields if provided
	if req.Title != nil {
		existingTodo.Title = *req.Title
	}
	if req.Description != nil {
		existingTodo.Description = *req.Description
	}
	if req.Status != nil {
		existingTodo.Status = *req.Status
	}

	updatedTodo, err := uc.todoRepo.Update(existingTodo)
	if err != nil {
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}

	return updatedTodo, nil
}

func (uc *todoUseCase) Delete(todoID, userID int, isAdmin bool) error {
	var err error

	if isAdmin {
		// Admin can delete any todo
		err = uc.todoRepo.Delete(todoID)
	} else {
		// User can only delete their own todo
		// First check if todo belongs to user
		_, err = uc.todoRepo.GetByIDAndUserID(todoID, userID)
		if err != nil {
			return fmt.Errorf("todo not found or access denied: %w", err)
		}

		err = uc.todoRepo.Delete(todoID)
	}

	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	return nil
}
