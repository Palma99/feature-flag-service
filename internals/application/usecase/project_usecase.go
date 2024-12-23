package usecase

import (
	"errors"
	"fmt"

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

func (i *ProjectInteractor) CreateProject(name string, userId int) (*int64, error) {

	_, err := i.keyService.GeneratePublicKey()
	if err != nil {
		return nil, errors.New("error during creation of project")
	}

	if name == "" {
		return nil, errors.New("error during creation of project, name is required")
	}

	project := &repository.CreateProjectDTO{
		Name:    name,
		OwnerId: userId,
	}

	created, err := i.projectRepository.CreateProject(project)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("error during creation of project")
	}

	return &created.ID, nil
}

func (i *ProjectInteractor) GetProjectsByUserId(userId int) ([]entity.Project, error) {
	return i.projectRepository.GetProjectsByUserId(userId)
}

func (i *ProjectInteractor) GetProjectDetails(userId int, projectId int64) (*entity.ProjectWithMembers, error) {
	projectWithMembers, err := i.projectRepository.GetProjectDetails(projectId)
	if err != nil || !projectWithMembers.HasMember(userId) {
		fmt.Println(err)
		return nil, errors.New("project not found")
	}

	project, err := i.projectRepository.GetProjectDetails(projectId)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("project not found")
	}

	return project, nil
}

func (i *ProjectInteractor) GetUserPermissionsOnThisProject(userId int, projectId int64) ([]string, error) {
	projectWithMembers, err := i.projectRepository.GetProjectDetails(projectId)
	if err != nil || !projectWithMembers.HasMember(userId) {
		fmt.Println(err)
		return nil, errors.New("project not found")
	}

	return projectWithMembers.GetUserPermissions(userId), nil
}
