package repository

import "database/sql"

type Repositories struct {
	User                 UserRepository
	Auth                 AuthRepository
	VideoServiceInstance VideoServiceInstanceRepository
	UserClient           UserClientRepository
}

func New() *Repositories {
	return &Repositories{
		User: NewUserRepository(),
		Auth: NewAuthRepository(),
	}
}

func NewWithDB(db *sql.DB) *Repositories {
	return &Repositories{
		User:                 NewUserRepository(),
		Auth:                 NewAuthRepository(),
		VideoServiceInstance: NewVideoServiceInstanceRepository(db),
		UserClient:           NewUserClientRepository(db),
	}
}
