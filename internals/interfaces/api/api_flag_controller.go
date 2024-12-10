package api_controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Palma99/feature-flag-service/internals/application/usecase"
	context_keys "github.com/Palma99/feature-flag-service/internals/infrastructure/context"
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
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))

		return
	}

	w.WriteHeader(http.StatusCreated)
}
