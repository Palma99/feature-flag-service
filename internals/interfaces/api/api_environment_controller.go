package api_controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Palma99/feature-flag-service/internals/application/usecase"
	context_keys "github.com/Palma99/feature-flag-service/internals/infrastructure/context"
	"github.com/go-chi/chi/v5"
)

type ApiEnvironmentController struct {
	environmentInteractor *usecase.EnvironmentInteractor
}

func NewApiEnvironmentController(environmentInteractor *usecase.EnvironmentInteractor) *ApiEnvironmentController {
	return &ApiEnvironmentController{
		environmentInteractor,
	}
}

func (environmentController *ApiEnvironmentController) GetEnvironment(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(context_keys.UserIDKey).(int)
	environmentId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	environment, err := environmentController.environmentInteractor.GetEnvironmentDetails(userId, environmentId)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonResponse := map[string]interface{}{
		"environment": environment,
	}

	json.NewEncoder(w).Encode(jsonResponse)
	w.WriteHeader(http.StatusOK)
}
