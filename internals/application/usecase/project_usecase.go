package usecase

import (
	"errors"

	"github.com/Palma99/feature-flag-service/internals/application/services"
	entity "github.com/Palma99/feature-flag-service/internals/domain/entity"
	repository "github.com/Palma99/feature-flag-service/internals/domain/repository"
)

type ProjectInteractor struct {
	projectRepository repository.ProjectRepository
	keyService        *services.KeyService
}

func NewProjectInteractor(
	projectRepository repository.ProjectRepository,
	keyService *services.KeyService,
) *ProjectInteractor {
	return &ProjectInteractor{
		projectRepository,
		keyService,
	}
}

func (i *ProjectInteractor) CreateEnvironment(envName string) error {
	// secret, err := i.keyService.GenerateSecretKey()
	// if err != nil {
	// 	return errors.New("error during creation of environment")
	// }

	// pk, err := i.keyService.GeneratePublicKey()
	// if err != nil {
	// 	return errors.New("error during creation of environment")
	// }

	// environment := &entity.Environment{
	// 	Name:       "My environment",
	// 	PublicKey:  pk,
	// 	PrivateKey: secret,
	// 	ProjectID:  1,
	// }

	return nil
}

func (i *ProjectInteractor) CreateProject(name string, privateKey string) (*entity.Environment, error) {
	if i.keyService.IsPublicKey(privateKey) {
		return nil, errors.New("unauthorized")
	}

	return nil, errors.New("not implemented")
}
