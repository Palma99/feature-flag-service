package usecase

import (
	"errors"
	"fmt"

	"github.com/Palma99/feature-flag-service/internals/application/services"
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

func (i *ProjectInteractor) CreateProject(name string, userId int) error {

	_, err := i.keyService.GeneratePublicKey()
	if err != nil {
		return errors.New("error during creation of project")
	}

	project := &repository.CreateProjectDTO{
		Name:    name,
		OwnerId: userId,
	}

	_, err = i.projectRepository.CreateProject(project)
	if err != nil {
		fmt.Println(err)
		return errors.New("error during creation of project")
	}

	return nil
}
