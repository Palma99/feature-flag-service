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

func (i *ProjectInteractor) CreateProject(name string, privateKey string) (*entity.Environment, error) {
	if i.keyService.IsPublicKey(privateKey) {
		return nil, errors.New("unauthorized")
	}

	return nil, errors.New("not implemented")
}
