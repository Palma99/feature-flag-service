package usecase

import (
	"errors"

	"github.com/Palma99/feature-flag-service/internals/application/services"
	entity "github.com/Palma99/feature-flag-service/internals/domain/entity"
	repository "github.com/Palma99/feature-flag-service/internals/domain/repository"
)

type EnvironmentInteractor struct {
	environmentRepository repository.EnvironmentRepository
	keyService            *services.KeyService
}

func NewEnvironmentInteractor(
	environmentRepository repository.EnvironmentRepository,
	keyService *services.KeyService,
) *EnvironmentInteractor {
	return &EnvironmentInteractor{
		environmentRepository,
		keyService,
	}
}

func (i *EnvironmentInteractor) CreateEnvironment(envName string, projectId string) error {

	if (envName == "") || (projectId == "") {
		return errors.New("error during creation of environment, name and project id are required")
	}

	secret, err := i.keyService.GenerateSecretKey()
	if err != nil {
		return errors.New("error during creation of environment")
	}

	pk, err := i.keyService.GeneratePublicKey()
	if err != nil {
		return errors.New("error during creation of environment")
	}

	environment := &entity.Environment{
		Name:       envName,
		PublicKey:  pk,
		PrivateKey: secret,
		ProjectID:  projectId,
	}

	err = i.environmentRepository.CreateEnvironment(environment)
	if err != nil {
		return errors.New("error during creation of environment")
	}

	return nil
}

func (i *EnvironmentInteractor) GetEnvironmentByKey(key string) (*entity.Environment, error) {
	if !i.keyService.IsPublicKey(key) {
		return nil, errors.New("secret key is not supported yet")
	}

	return i.environmentRepository.GetEnvironmentByPublicKey(key)
}
