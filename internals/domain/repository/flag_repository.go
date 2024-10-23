package domain

import entity "github.com/Palma99/feature-flag-service/internals/domain/entity"

type FlagRepository interface {
	GetAllFlagsByEnvironmentID(environmentID int) ([]entity.Flag, error)
	GetFlagInEnvironmentById(environmentID, flagId int) (*entity.Flag, error)
	UpdateEnvironmentFlagValue(environmentId int, flag *entity.Flag) error
}
