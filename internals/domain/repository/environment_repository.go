package domain

import entity "github.com/Palma99/feature-flag-service/internals/domain/entity"

type EnvironmentRepository interface {
	// todo: spostare nel repository di flag
	GetEnvironmentActiveFlagsByPublicKey(publicKey string) ([]string, error)
	CreateEnvironment(env *entity.Environment) error
	GetEnvironmentDetails(id int64) (*entity.EnvironmentWithFlags, error)
}
