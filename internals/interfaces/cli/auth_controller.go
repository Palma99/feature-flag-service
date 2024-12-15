package interfaces

import "github.com/Palma99/feature-flag-service/internals/application/usecase"

type AuthController struct {
	authInteractor *usecase.AuthInteractor
}

func NewUserController(authInteractor *usecase.AuthInteractor) *AuthController {
	return &AuthController{
		authInteractor,
	}
}

func (authController *AuthController) CreateUser(userDTO usecase.CreateUserDTO) error {
	err := authController.authInteractor.CreateUser(userDTO)

	return err
}

func (authController *AuthController) GetToken(loginDTO usecase.UserLoginDTO) (*string, error) {
	return authController.authInteractor.GetToken(loginDTO)
}
