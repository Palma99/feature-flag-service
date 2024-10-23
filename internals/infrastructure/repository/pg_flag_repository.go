package infrastructure

import (
	"database/sql"

	entity "github.com/Palma99/feature-flag-service/internals/domain/entity"
)

type PgFlagRepository struct {
	db *sql.DB
}

func NewPgFlagRepository(db *sql.DB) *PgFlagRepository {
	return &PgFlagRepository{
		db,
	}
}

func (r *PgFlagRepository) UpdateEnvironmentFlagValue(environmentId int, flag *entity.Flag) error {

	_, err := r.db.Exec(`
		UPDATE flag_environment
		SET enabled = $1
		WHERE environment = $2 AND flag = $3
	`, flag.Enabled, environmentId, flag.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *PgFlagRepository) GetFlagInEnvironmentById(environmentID, flagId int) (*entity.Flag, error) {

	flagRow := r.db.QueryRow(`
		SELECT
			f.id as flag_id,
			name,
			enabled
			FROM flag_environment fe
			LEFT JOIN flag f on f.id = fe.flag
			WHERE f.id = $1 AND environment = $2
	`, flagId, environmentID)

	if flagRow.Err() != nil {
		return nil, flagRow.Err()
	}

	flag := &entity.Flag{}

	if err := flagRow.Scan(
		&flag.ID,
		&flag.Name,
		&flag.Enabled,
	); err != nil {
		return nil, err
	}

	return flag, nil
}

func (r *PgFlagRepository) GetAllFlagsByEnvironmentID(environmentID int) ([]entity.Flag, error) {

	rows, err := r.db.Query(`
		SELECT
			e.id as env_id,
			f.id as flag_id,
			e.name as env_name,
			f.name as flag_name,
			fe.enabled
			FROM flag_environment fe 
			LEFT JOIN flag f on f.id = fe.flag 
			LEFT JOIN environment e on e.id = fe.environment
			WHERE fe.environment = $1;
	`, environmentID)

	if err != nil {
		return nil, err
	}

	flags := []entity.Flag{}

	for rows.Next() {
		flag := entity.Flag{}
		if err := rows.Scan(
			&flag.Environment.ID,
			&flag.ID,
			&flag.Environment.Name,
			&flag.Name,
			&flag.Enabled,
		); err != nil {
			return nil, err
		}
		flags = append(flags, flag)
	}

	return flags, nil
}
