package interfaces

import "github.com/Palma99/feature-flag-service/internals/application/usecase"

type ProjectController struct {
	projectInteractor *usecase.ProjectInteractor
}

func NewProjectController(projectInteractor *usecase.ProjectInteractor) *ProjectController {
	return &ProjectController{
		projectInteractor,
	}
}

func (projectController *ProjectController) CreateProject(name string, userId int) error {
	return projectController.projectInteractor.CreateProject(name, userId)
}
