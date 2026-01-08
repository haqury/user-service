package service

import (
	"context"
	"github.com/haqury/user-service/internal/repository"
)

type AuthService interface {
	Login(ctx context.Context, username, password string) (string, interface{}, error)
	ValidateToken(ctx context.Context, token string) (bool, interface{}, error)
	Logout(ctx context.Context, token string) error
}

type authService struct {
	authRepo repository.AuthRepository
	userRepo repository.UserRepository
}

func NewAuthService(authRepo repository.AuthRepository) AuthService {
	return &authService{authRepo: authRepo}
}

func (s *authService) Login(ctx context.Context, username, password string) (string, interface{}, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", nil, err
	}

	token, err := s.authRepo.CreateToken(ctx, "user-id")
	return token, user, err
}

func (s *authService) ValidateToken(ctx context.Context, token string) (bool, interface{}, error) {
	userID, err := s.authRepo.ValidateToken(ctx, token)
	if err != nil {
		return false, nil, err
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	return true, user, err
}

func (s *authService) Logout(ctx context.Context, token string) error {
	return s.authRepo.RevokeToken(ctx, token)
}
