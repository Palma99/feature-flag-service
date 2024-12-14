package api_controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Palma99/feature-flag-service/internals/application/usecase"
	context_keys "github.com/Palma99/feature-flag-service/internals/infrastructure/context"
)

type PublicServiceController struct {
	environmentInteractor *usecase.EnvironmentInteractor
}

func NewPublicServiceController(environmentService *usecase.EnvironmentInteractor) *PublicServiceController {
	return &PublicServiceController{
		environmentService,
	}
}

func (c *PublicServiceController) GetFlagsByPublicKey(w http.ResponseWriter, r *http.Request) {
	publicKey := r.Context().Value(context_keys.PublicKeyKey).(string)

	flags, err := c.environmentInteractor.GetEnvironmentFlagsByPublicKey(publicKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonResponse := map[string]interface{}{
		"activeFlags": flags,
	}

	json.NewEncoder(w).Encode(jsonResponse)
	w.WriteHeader(http.StatusOK)
}
