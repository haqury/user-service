package service

import "github.com/haqury/user-service/internal/repository"

type Services struct {
	User    UserService
	Auth    AuthService
	Routing RoutingService
}

func New(repos *repository.Repositories) *Services {
	return &Services{
		User:    NewUserService(repos.User, repos.UserClient),
		Auth:    NewAuthService(repos.Auth),
		Routing: NewRoutingService(repos),
	}
}
