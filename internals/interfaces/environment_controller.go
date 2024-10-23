package interfaces

import (
	"github.com/Palma99/feature-flag-service/internals/application/usecase"
)

type EnvironmentController struct {
	environmentInteractor *usecase.EnvironmentInteractor
}

func NewEnvironmentInteractor(environmentInteractor *usecase.EnvironmentInteractor) *EnvironmentController {
	return &EnvironmentController{
		environmentInteractor: environmentInteractor,
	}
}

func (ec *EnvironmentController) CreateEnvironment(name, projectId string) error {
	return ec.environmentInteractor.CreateEnvironment(name, projectId)
}
