package repository

import (
	"context"
)

type AuthRepository interface {
	CreateToken(ctx context.Context, userID string) (string, error)
	ValidateToken(ctx context.Context, token string) (string, error)
	RevokeToken(ctx context.Context, token string) error
}

type authRepository struct{}

func NewAuthRepository() AuthRepository {
	return &authRepository{}
}

func (r *authRepository) CreateToken(ctx context.Context, userID string) (string, error) {
	return "dummy-token-" + userID, nil
}

func (r *authRepository) ValidateToken(ctx context.Context, token string) (string, error) {
	return "user-id-from-token", nil
}

func (r *authRepository) RevokeToken(ctx context.Context, token string) error {
	return nil
}
