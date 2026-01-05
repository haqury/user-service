package repository

import (
	"context"
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (interface{}, error)
	GetByUsername(ctx context.Context, username string) (interface{}, error)
	Create(ctx context.Context, user interface{}) error
	Update(ctx context.Context, id string, user interface{}) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int, filter string) ([]interface{}, int, error)
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) GetByID(ctx context.Context, id string) (interface{}, error) {
	return map[string]interface{}{
		"id":       id,
		"username": "testuser",
		"email":    "test@example.com",
		"status":   "active",
	}, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (interface{}, error) {
	return map[string]interface{}{
		"id":       "123",
		"username": username,
		"email":    "user@example.com",
		"status":   "active",
	}, nil
}

func (r *userRepository) Create(ctx context.Context, user interface{}) error {
	return nil
}

func (r *userRepository) Update(ctx context.Context, id string, user interface{}) error {
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *userRepository) List(ctx context.Context, page, limit int, filter string) ([]interface{}, int, error) {
	return []interface{}{
		map[string]interface{}{
			"id":       "1",
			"username": "user1",
			"email":    "user1@example.com",
		},
		map[string]interface{}{
			"id":       "2",
			"username": "user2",
			"email":    "user2@example.com",
		},
	}, 2, nil
}
