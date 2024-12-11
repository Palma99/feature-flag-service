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

type CreateEnvironmentDTO struct {
	Name      string `json:"name"`
	ProjectId int64  `json:"projectId"`
}

func (environmentController *ApiEnvironmentController) CreateEnvironment(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(context_keys.UserIDKey).(int)

	var dto CreateEnvironmentDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = environmentController.environmentInteractor.CreateEnvironment(dto.Name, dto.ProjectId, userId)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
}
