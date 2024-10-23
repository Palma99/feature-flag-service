package usecase

import (
	"github.com/Palma99/feature-flag-service/internals/application/services"
	entity "github.com/Palma99/feature-flag-service/internals/domain/entity"
	repository "github.com/Palma99/feature-flag-service/internals/domain/repository"
)

type CreateUserDTO struct {
	Email    string
	Nickname string
	Password string
}

type UserInteractor struct {
	userRepository  repository.UserRepository
	passwordService services.PasswordService
}

func NewUserInteractor(
	userRepository repository.UserRepository,
	passwordService services.PasswordService,
) *UserInteractor {
	return &UserInteractor{
		userRepository,
		passwordService,
	}
}

func (u *UserInteractor) CreateUser(createUserDTO CreateUserDTO) error {

	hashedPassword, err := u.passwordService.HashPassword(createUserDTO.Password)
	if err != nil {
		return err
	}

	user := entity.User{
		Email:    createUserDTO.Email,
		Password: hashedPassword,
		Nickname: createUserDTO.Email,
	}

	return u.userRepository.CreateUser(&user)
}
