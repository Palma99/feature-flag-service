package api_controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Palma99/feature-flag-service/internals/application/usecase"
)

type ApiAuthController struct {
	authInteractor *usecase.AuthInteractor
}

func NewApiUserController(authInteractor *usecase.AuthInteractor) *ApiAuthController {
	return &ApiAuthController{
		authInteractor,
	}
}

func (authController *ApiAuthController) CreateUser(userDTO usecase.CreateUserDTO) error {
	err := authController.authInteractor.CreateUser(userDTO)

	return err
}

func (authController *ApiAuthController) GetToken(w http.ResponseWriter, r *http.Request) {
	var loginDTO usecase.UserLoginDTO

	err := json.NewDecoder(r.Body).Decode(&loginDTO)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := authController.authInteractor.GetToken(loginDTO)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)

	jsonResponse := map[string]string{
		"token": *token,
	}

	json.NewEncoder(w).Encode(jsonResponse)
}
