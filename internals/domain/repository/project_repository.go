package domain

import (
	entity "github.com/Palma99/feature-flag-service/internals/domain/entity"
)

type ProjectRepository interface {
	CreateProject(project *entity.Project) error
}
