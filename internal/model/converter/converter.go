package converter

import (
	"github.com/islamyakin/otel-propagation-monorepo/internal/entity"
	"github.com/islamyakin/otel-propagation-monorepo/internal/model"
)

func UserModelToEntity(m *model.UserModel) *entity.User {
	if m == nil {
		return nil
	}
	return &entity.User{
		ID:        m.ID,
		Username:  m.Username,
		Password:  m.Password,
		Role:      entity.Role(m.Role),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func UserEntityToModel(e *entity.User) *model.UserModel {
	if e == nil {
		return nil
	}
	return &model.UserModel{
		ID:        e.ID,
		Username:  e.Username,
		Password:  e.Password,
		Role:      string(e.Role),
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func TodoModelToEntity(m *model.TodoModel) *entity.Todo {
	if m == nil {
		return nil
	}
	return &entity.Todo{
		ID:          m.ID,
		UserID:      m.UserID,
		Title:       m.Title,
		Description: m.Description,
		Status:      entity.TodoStatus(m.Status),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func TodoEntityToModel(e *entity.Todo) *model.TodoModel {
	if e == nil {
		return nil
	}
	return &model.TodoModel{
		ID:          e.ID,
		UserID:      e.UserID,
		Title:       e.Title,
		Description: e.Description,
		Status:      string(e.Status),
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func TodoModelsToEntities(models []*model.TodoModel) []*entity.Todo {
	entities := make([]*entity.Todo, len(models))
	for i, m := range models {
		entities[i] = TodoModelToEntity(m)
	}
	return entities
}

func UserModelsToEntities(models []*model.UserModel) []*entity.User {
	entities := make([]*entity.User, len(models))
	for i, m := range models {
		entities[i] = UserModelToEntity(m)
	}
	return entities
}
