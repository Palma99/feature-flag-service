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

func (r *PgFlagRepository) CreateFlag(flag *domain.Flag) error {
	_, err := r.db.Exec(`INSERT INTO flag (name, project_id) VALUES ($1, $2)`,
		flag.Name, flag.ProjectID,
	)
	if err != nil {
		return err
	}

	return nil
}
