package usecase

import (
	"errors"
	"fmt"

	"github.com/Palma99/feature-flag-service/internals/application/services"
	entity "github.com/Palma99/feature-flag-service/internals/domain/entity"
	repository "github.com/Palma99/feature-flag-service/internals/domain/repository"
)

type EnvironmentInteractor struct {
	environmentRepository repository.EnvironmentRepository
	projectRepository     repository.ProjectRepository
	keyService            *services.KeyService
}

func NewEnvironmentInteractor(
	environmentRepository repository.EnvironmentRepository,
	projectRepository repository.ProjectRepository,
	keyService *services.KeyService,
) *EnvironmentInteractor {
	return &EnvironmentInteractor{
		environmentRepository,
		projectRepository,
		keyService,
	}
}

func (i *EnvironmentInteractor) CreateEnvironment(envName string, projectId int64, userId int) error {

	if (envName == "") || (projectId == 0) {
		return errors.New("error during creation of environment, name and project id are required")
	}

	pk, err := i.keyService.GeneratePublicKey()
	if err != nil {
		return errors.New("error during creation of environment")
	}

	loggedUser := &entity.LoggedUser{
		ID: userId,
	}

	project, err := i.projectRepository.GetUserProjectByProjectId(userId, projectId)
	if err != nil {
		fmt.Println(err)
		return errors.New("error during creation of environment")
	}

	if !loggedUser.CanCreateProjectEnvironment(*project) {
		return errors.New("user is not allowed to create this environment")
	}

	environment := &entity.Environment{
		Name:      envName,
		PublicKey: pk,
		ProjectID: projectId,
	}

	if err := project.CanCreateEnvironment(*environment); err != nil {
		return err
	}

	err = i.environmentRepository.CreateEnvironment(environment)
	if err != nil {
		fmt.Println(err)
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

func (i *EnvironmentInteractor) GetEnvironmentDetails(userId int, environmentId int64) (*entity.EnvironmentWithFlags, error) {

	environmentDetails, err := i.environmentRepository.GetEnvironmentDetails(environmentId)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("environment not found")
	}

	userProjectRelationId, err := i.projectRepository.GetUserProjectRelation(userId, environmentDetails.ProjectID)

	if err != nil || userProjectRelationId == nil {
		fmt.Println(err)
		return nil, errors.New("project not found")
	}

	return environmentDetails, nil
}
