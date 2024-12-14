package api_controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Palma99/feature-flag-service/internals/application/usecase"
	domain "github.com/Palma99/feature-flag-service/internals/domain/entity"
	context_keys "github.com/Palma99/feature-flag-service/internals/infrastructure/context"
	"github.com/go-chi/chi/v5"
)

type ApiFlagController struct {
	flagInteractor *usecase.FlagInteractor
}

func NewApiFlagController(flagInteractor *usecase.FlagInteractor) *ApiFlagController {
	return &ApiFlagController{
		flagInteractor,
	}
}

type CreateFlagDTO struct {
	FlagName  string `json:"flagName"`
	ProjectId int64  `json:"projectId"`
}

func (flagController *ApiFlagController) GetProjectFlags(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(context_keys.UserIDKey).(int)
	projectId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	flags, err := flagController.flagInteractor.GetProjectFlags(userId, projectId)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonResponse := map[string]interface{}{
		"projectFlags": flags,
	}

	json.NewEncoder(w).Encode(jsonResponse)
	w.WriteHeader(http.StatusOK)
}

func (flagController *ApiFlagController) UpdateFlagEnvironment(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(context_keys.UserIDKey).(int)
	environmentId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	flagsToUpdate := []domain.Flag{}
	err = json.NewDecoder(r.Body).Decode(&flagsToUpdate)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = flagController.flagInteractor.UpdateFlagEnvironment(int(environmentId), userId, flagsToUpdate)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (flagController *ApiFlagController) CreateFlag(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(context_keys.UserIDKey).(int)

	createFlagDTO := CreateFlagDTO{}
	err := json.NewDecoder(r.Body).Decode(&createFlagDTO)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = flagController.flagInteractor.CreateFlag(createFlagDTO.ProjectId, userId, createFlagDTO.FlagName)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (flagController *ApiFlagController) DeleteFlag(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(context_keys.UserIDKey).(int)
	flagId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	err = flagController.flagInteractor.DeleteFlag(userId, int(flagId))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	w.WriteHeader(http.StatusOK)
}
