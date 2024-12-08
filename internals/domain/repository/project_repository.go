package domain

import (
	entity "github.com/Palma99/feature-flag-service/internals/domain/entity"
)

type ProjectRepository interface {
	CreateProject(project *CreateProjectDTO) (*entity.Project, error)
	GetProjectsByUserId(userId int) ([]entity.Project, error)
	GetUserProjectByProjectId(userId int, projectId int64) (*entity.Project, error)
}

type CreateProjectDTO struct {
	Name    string
	OwnerId int
}
