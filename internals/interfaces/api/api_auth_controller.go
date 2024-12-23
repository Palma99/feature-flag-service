package api_controllers

import (
	"encoding/json"
	"fmt"
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

func (authController *ApiAuthController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var signupDTO usecase.CreateUserDTO

	err := json.NewDecoder(r.Body).Decode(&signupDTO)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = authController.authInteractor.CreateUser(signupDTO)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse := map[string]string{
			"error": err.Error(),
		}

		json.NewEncoder(w).Encode(jsonResponse)
		return
	}

	w.WriteHeader(http.StatusCreated)
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
