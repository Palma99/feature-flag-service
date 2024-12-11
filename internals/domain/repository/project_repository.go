package domain

import (
	entity "github.com/Palma99/feature-flag-service/internals/domain/entity"
)

type ProjectRepository interface {
	CreateProject(project *CreateProjectDTO) (*entity.Project, error)
	GetProjectsByUserId(userId int) ([]entity.Project, error)
	GetUserProjectRelation(userId int, projectId int64) (*int64, error)
	GetProjectDetails(projectId int64) (*entity.ProjectWithMembers, error)
}

type CreateProjectDTO struct {
	Name    string `json:"projectName"`
	OwnerId int
}
