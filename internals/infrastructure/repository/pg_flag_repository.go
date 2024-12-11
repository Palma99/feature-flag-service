package infrastructure

import (
	"database/sql"

	domain "github.com/Palma99/feature-flag-service/internals/domain/entity"
)

type PgFlagRepository struct {
	db *sql.DB
}

func NewPgFlagRepository(db *sql.DB) *PgFlagRepository {
	return &PgFlagRepository{
		db,
	}
}

func (r *PgFlagRepository) UpdateFlagEnvironment(environmentId int, flagsToUpdate []domain.Flag) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	for _, flag := range flagsToUpdate {
		if err := upsertFlagEnvironment(tx, environmentId, flag); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func upsertFlagEnvironment(tx *sql.Tx, environmentId int, flag domain.Flag) error {
	_, err := tx.Exec(`INSERT INTO flag_environment (flag, environment, enabled)
		VALUES ($1, $2, $3)
		ON CONFLICT (flag, environment) DO UPDATE
		SET enabled = EXCLUDED.enabled;`, flag.ID, environmentId, flag.Enabled)

	return err
}

func (r *PgFlagRepository) CreateFlag(flag *domain.Flag) error {
	_, err := r.db.Exec(`INSERT INTO flag (name, project_id) VALUES ($1, $2)`,
		flag.Name, flag.ProjectID,
	)
	if err != nil {
		return err
	}

	return nil
}
