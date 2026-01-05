package service

import (
	"context"
	"user-service/internal/repository"
)

type UserService interface {
	GetUser(ctx context.Context, id string) (interface{}, error)
	GetUserByUsername(ctx context.Context, username string) (interface{}, error)
	CreateUser(ctx context.Context, user interface{}) error
	UpdateUser(ctx context.Context, id string, user interface{}) error
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, page, limit int, filter string) ([]interface{}, int, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUser(ctx context.Context, id string) (interface{}, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (interface{}, error) {
	return s.repo.GetByUsername(ctx, username)
}

func (s *userService) CreateUser(ctx context.Context, user interface{}) error {
	return s.repo.Create(ctx, user)
}

func (s *userService) UpdateUser(ctx context.Context, id string, user interface{}) error {
	return s.repo.Update(ctx, id, user)
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *userService) ListUsers(ctx context.Context, page, limit int, filter string) ([]interface{}, int, error) {
	return s.repo.List(ctx, page, limit, filter)
}
