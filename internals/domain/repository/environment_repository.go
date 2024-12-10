package domain

import entity "github.com/Palma99/feature-flag-service/internals/domain/entity"

type EnvironmentRepository interface {
	GetEnvironmentByPublicKey(publicKey string) (*entity.Environment, error)
	GetEnvironmentBySecretKey(secretKey string) (*entity.Environment, error)
	CreateEnvironment(env *entity.Environment) error
	GetEnvironmentDetails(id int64) (*entity.EnvironmentWithFlags, error)
}
