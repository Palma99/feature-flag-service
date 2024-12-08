package api_controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Palma99/feature-flag-service/internals/application/usecase"
	domain "github.com/Palma99/feature-flag-service/internals/domain/repository"
	context_keys "github.com/Palma99/feature-flag-service/internals/infrastructure/context"
)

type ApiProjectController struct {
	projectInteractor *usecase.ProjectInteractor
}

func NewApiProjectController(projectInteractor *usecase.ProjectInteractor) *ApiProjectController {
	return &ApiProjectController{
		projectInteractor,
	}
}

func (projectController *ApiProjectController) CreateProject(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(context_keys.UserIDKey).(int)

	var dto domain.CreateProjectDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dto.OwnerId = userId

	createdId, err := projectController.projectInteractor.CreateProject(dto.Name, dto.OwnerId)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonResponse := map[string]interface{}{
		"projectId": *createdId,
	}
	json.NewEncoder(w).Encode(jsonResponse)
	w.WriteHeader(http.StatusCreated)
}

func (projectController *ApiProjectController) GetProjects(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(context_keys.UserIDKey).(int)

	projects, err := projectController.projectInteractor.GetProjectsByUserId(userId)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	jsonResponse := map[string]interface{}{
		"projects": projects,
	}

	json.NewEncoder(w).Encode(jsonResponse)
}
