package infrastructure

import (
	"database/sql"
	"encoding/json"

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
		INSERT INTO environment (name, public_key, project_id) VALUES ($1, $2, $3)
	`, env.Name, env.PublicKey, env.ProjectID)

	if err != nil {
		return err
	}

	return nil
}

func (r *PgEnvironmentRepository) getProjectFlags(projectId int64) ([]entity.Flag, error) {
	rows, err := r.db.Query(`
		SELECT id, name, project_id FROM flag WHERE project_id = $1
	`, projectId)

	if err != nil {
		return nil, err
	}

	var flags []entity.Flag

	for rows.Next() {
		flag := entity.Flag{
			Enabled: false,
		}
		if err := rows.Scan(&flag.ID, &flag.Name, &flag.ProjectID); err != nil {
			return nil, err
		}
		flags = append(flags, flag)
	}

	return flags, nil
}

func (r *PgEnvironmentRepository) GetEnvironmentDetails(environmentId int64) (*entity.EnvironmentWithFlags, error) {

	row := r.db.QueryRow(`
		SELECT 
			e.ID AS EnvironmentID,
			e.Name AS EnvironmentName,
			e.project_id ,
			e.public_key,
			COALESCE(
				JSON_AGG(
					JSON_BUILD_OBJECT(
						'id', f.ID,
						'name', f.Name,
            'project_id', f.project_id,
						'enabled', COALESCE(fe.Enabled, FALSE)
					)
				) FILTER (WHERE f.ID IS NOT NULL),
				'[]' -- Array vuoto se non ci sono flag
			) AS Flags
		FROM 
				Environment e
		LEFT JOIN 
				Flag_Environment fe ON e.ID = fe.environment 
		LEFT JOIN 
				Flag f ON f.ID = fe.flag
		WHERE 
				e.ID = $1
		GROUP BY 
				e.ID, e.Name, e.project_id ;
	`, environmentId)

	var envId int
	var envName string
	var projectId int64
	var publicKey string
	var flagsJSON string

	if err := row.Scan(&envId, &envName, &projectId, &publicKey, &flagsJSON); err != nil {
		return nil, err
	}

	env := &entity.EnvironmentWithFlags{
		Environment: entity.Environment{
			ID:        envId,
			Name:      envName,
			PublicKey: publicKey,
			ProjectID: projectId,
		},
		Flags: []entity.Flag{},
	}

	if err := json.Unmarshal([]byte(flagsJSON), &env.Flags); err != nil {
		return nil, err
	}

	// merge with project flags
	projectFlags, err := r.getProjectFlags(projectId)
	if err != nil {
		return nil, err
	}

	for _, projectFlag := range projectFlags {
		flagHasConfigurationInEnvironment := false
		for _, envFlag := range env.Flags {
			if envFlag.ID == projectFlag.ID {
				flagHasConfigurationInEnvironment = true
				break
			}
		}
		if !flagHasConfigurationInEnvironment {
			env.Flags = append(env.Flags, projectFlag)
		}
	}

	return env, nil
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
