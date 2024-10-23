package interfaces

import "github.com/Palma99/feature-flag-service/internals/application/usecase"

type UserController struct {
	userInteractor *usecase.UserInteractor
}

func NewUserController(userInteractor *usecase.UserInteractor) *UserController {
	return &UserController{
		userInteractor,
	}
}

func (userController *UserController) CreateUser(userDTO usecase.CreateUserDTO) error {
	err := userController.userInteractor.CreateUser(userDTO)

	return err
}
