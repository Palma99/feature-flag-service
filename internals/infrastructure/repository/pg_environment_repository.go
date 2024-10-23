package infrastructure

import (
	"database/sql"

	entity "github.com/Palma99/feature-flag-service/internals/domain/entity"
)

type PgEnvironmentRepository struct {
	db *sql.DB
}

func NewPgEnvironmentRepository(db *sql.DB) *PgEnvironmentRepository {
	return &PgEnvironmentRepository{
		db,
	}
}

func (r *PgEnvironmentRepository) GetEnvironmentByPublicKey(key string) (*entity.Environment, error) {
	envRow := r.db.QueryRow(`
		SELECT id, name FROM environment WHERE public_key = $1
	`, key)

	if envRow.Err() != nil {
		return nil, envRow.Err()
	}

	env := entity.Environment{}

	if err := envRow.Scan(&env.ID, &env.Name); err != nil {
		return nil, err
	}

	return &env, nil
}

func (r *PgEnvironmentRepository) CreateEnvironment(env *entity.Environment) error {
	_, err := r.db.Exec(`
		INSERT INTO environment (name, public_key, private_key, project_id) VALUES ($1, $2, $3, $4)
	`, env.Name, env.PublicKey, env.PrivateKey, env.ProjectID)

	if err != nil {
		return err
	}

	return nil
}

func (r *PgEnvironmentRepository) GetEnvironmentByName(name string) (*entity.Environment, error) {

	return nil, nil
}

func (r *PgEnvironmentRepository) GetEnvironmentBySecretKey(key string) (*entity.Environment, error) {
	envRow := r.db.QueryRow(`
		SELECT id, name FROM environment WHERE private_key = $1
	`, key)

	if envRow.Err() != nil {
		return nil, envRow.Err()
	}

	env := entity.Environment{}

	if err := envRow.Scan(&env.ID, &env.Name); err != nil {
		return nil, err
	}

	return &env, nil
}
