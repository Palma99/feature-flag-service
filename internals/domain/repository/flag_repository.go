package domain

import entity "github.com/Palma99/feature-flag-service/internals/domain/entity"

type FlagRepository interface {
	CreateFlag(flag *entity.Flag) error
	UpdateFlagEnvironment(environmentId int, flag []entity.Flag) error
}
