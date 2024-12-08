package usecase

import (
	"errors"

	"github.com/Palma99/feature-flag-service/internals/application/services"
	entity "github.com/Palma99/feature-flag-service/internals/domain/entity"
	repository "github.com/Palma99/feature-flag-service/internals/domain/repository"
)

type CreateUserDTO struct {
	Email    string
	Nickname string
	Password string
}

type UserLoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthInteractor struct {
	userRepository  repository.UserRepository
	passwordService *services.PasswordService
	jwtService      *services.JWTService
}

func NewAuthInteractor(
	userRepository repository.UserRepository,
	passwordService *services.PasswordService,
	jwtService *services.JWTService,
) *AuthInteractor {
	return &AuthInteractor{
		userRepository,
		passwordService,
		jwtService,
	}
}

func (u *AuthInteractor) CreateUser(createUserDTO CreateUserDTO) error {
	hashedPassword, err := u.passwordService.HashPassword(createUserDTO.Password)
	if err != nil {
		return err
	}

	foundUser, _ := u.userRepository.GetUserByEmail(createUserDTO.Email)

	if foundUser != nil {
		return errors.New("user already exists")
	}

	user := entity.User{
		Email:    createUserDTO.Email,
		Password: hashedPassword,
		Nickname: createUserDTO.Nickname,
	}

	return u.userRepository.CreateUser(&user)
}

func (u *AuthInteractor) GetToken(userLoginDTO UserLoginDTO) (*string, error) {
	user, err := u.userRepository.GetUserByEmail(userLoginDTO.Email)

	if err != nil || user == nil {
		return nil, errors.New("invalid credentials")
	}

	if !u.passwordService.ArePasswordsEqual(userLoginDTO.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	token, err := u.jwtService.GenerateToken(user.ID)

	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (u *AuthInteractor) ValidateToken(token string) (int, error) {
	payload, err := u.jwtService.ValidateToken(token)

	if err != nil || payload == nil {
		return 0, err
	}

	return payload.UserID, nil
}
