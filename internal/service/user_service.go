package service

import (
	"context"

	"github.com/haqury/user-service/internal/repository"
)

type UserService interface {
	GetUser(ctx context.Context, id string) (interface{}, error)
	GetUserByUsername(ctx context.Context, username string) (interface{}, error)
	GetUserByClientID(ctx context.Context, clientID string) (*UserByClientIDResponse, error)
	CreateUser(ctx context.Context, user interface{}) error
	UpdateUser(ctx context.Context, id string, user interface{}) error
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, page, limit int, filter string) ([]interface{}, int, error)
}

type UserByClientIDResponse struct {
	UserID   string
	Username string
	IsActive bool
	Roles    []string
}

type userService struct {
	repo       repository.UserRepository
	userClient repository.UserClientRepository
}

func NewUserService(repo repository.UserRepository, userClient repository.UserClientRepository) UserService {
	return &userService{
		repo:       repo,
		userClient: userClient,
	}
}

func (s *userService) GetUser(ctx context.Context, id string) (interface{}, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (interface{}, error) {
	return s.repo.GetByUsername(ctx, username)
}

func (s *userService) GetUserByClientID(ctx context.Context, clientID string) (*UserByClientIDResponse, error) {
	// Получаем client запись
	userClient, err := s.userClient.GetByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	// Получаем пользователя
	user, err := s.repo.GetByID(ctx, userClient.UserID)
	if err != nil {
		return nil, err
	}

	// TODO: преобразовать user в нормальную структуру когда репозиторий будет готов
	// Сейчас это map[string]interface{}, нужно будет изменить на models.User
	userMap, ok := user.(map[string]interface{})
	if !ok {
		return nil, nil
	}

	return &UserByClientIDResponse{
		UserID:   userClient.UserID,
		Username: userMap["username"].(string),
		IsActive: true,             // TODO: get from user
		Roles:    []string{"user"}, // TODO: get from user
	}, nil
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
