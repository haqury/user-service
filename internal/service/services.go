package service

import "github.com/haqury/user-service/internal/repository"

type Services struct {
	User UserService
	Auth AuthService
}

func New(repos *repository.Repositories) *Services {
	return &Services{
		User: NewUserService(repos.User),
		Auth: NewAuthService(repos.Auth),
	}
}
