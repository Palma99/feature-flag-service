package domain

import (
	entity "github.com/Palma99/feature-flag-service/internals/domain/entity"
)

type UserRepository interface {
	CreateUser(user *entity.User) error
	GetUserByEmail(email string) (*entity.User, error)
}
