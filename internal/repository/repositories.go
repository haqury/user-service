package repository

type Repositories struct {
	User UserRepository
	Auth AuthRepository
}

func New() *Repositories {
	return &Repositories{
		User: NewUserRepository(),
		Auth: NewAuthRepository(),
	}
}
