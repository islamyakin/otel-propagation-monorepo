package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/islamyakin/otel-propagation-monorepo/internal/entity"
	"github.com/islamyakin/otel-propagation-monorepo/internal/model"
	"github.com/islamyakin/otel-propagation-monorepo/internal/model/converter"
)

type UserRepository interface {
	Create(user *entity.User) (*entity.User, error)
	GetByID(id int) (*entity.User, error)
	GetByUsername(username string) (*entity.User, error)
	GetAll() ([]*entity.User, error)
	Update(user *entity.User) (*entity.User, error)
	Delete(id int) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *entity.User) (*entity.User, error) {
	query := `
		INSERT INTO users (username, password, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, username, password, role, created_at, updated_at
	`

	now := time.Now()
	var userModel model.UserModel

	err := r.db.QueryRow(query, user.Username, user.Password, string(user.Role), now, now).
		Scan(&userModel.ID, &userModel.Username, &userModel.Password, &userModel.Role, &userModel.CreatedAt, &userModel.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return converter.UserModelToEntity(&userModel), nil
}

func (r *userRepository) GetByID(id int) (*entity.User, error) {
	query := `
		SELECT id, username, password, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var userModel model.UserModel
	err := r.db.QueryRow(query, id).
		Scan(&userModel.ID, &userModel.Username, &userModel.Password, &userModel.Role, &userModel.CreatedAt, &userModel.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return converter.UserModelToEntity(&userModel), nil
}

func (r *userRepository) GetByUsername(username string) (*entity.User, error) {
	query := `
		SELECT id, username, password, role, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	var userModel model.UserModel
	err := r.db.QueryRow(query, username).
		Scan(&userModel.ID, &userModel.Username, &userModel.Password, &userModel.Role, &userModel.CreatedAt, &userModel.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return converter.UserModelToEntity(&userModel), nil
}

func (r *userRepository) GetAll() ([]*entity.User, error) {
	query := `
		SELECT id, username, password, role, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	defer rows.Close()

	var userModels []*model.UserModel
	for rows.Next() {
		var userModel model.UserModel
		err := rows.Scan(&userModel.ID, &userModel.Username, &userModel.Password, &userModel.Role, &userModel.CreatedAt, &userModel.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		userModels = append(userModels, &userModel)
	}

	return converter.UserModelsToEntities(userModels), nil
}

func (r *userRepository) Update(user *entity.User) (*entity.User, error) {
	query := `
		UPDATE users
		SET username = $2, password = $3, role = $4, updated_at = $5
		WHERE id = $1
		RETURNING id, username, password, role, created_at, updated_at
	`

	now := time.Now()
	var userModel model.UserModel

	err := r.db.QueryRow(query, user.ID, user.Username, user.Password, string(user.Role), now).
		Scan(&userModel.ID, &userModel.Username, &userModel.Password, &userModel.Role, &userModel.CreatedAt, &userModel.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return converter.UserModelToEntity(&userModel), nil
}

func (r *userRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
