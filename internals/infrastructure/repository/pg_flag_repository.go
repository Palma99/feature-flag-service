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

func (r *PgFlagRepository) GetProjectFlags(projectId int64) ([]domain.FlagWithoutEnabled, error) {
	rows, err := r.db.Query(`SELECT id, name, project_id FROM flag WHERE project_id = $1`, projectId)
	if err != nil {
		return nil, err
	}

	flags := []domain.FlagWithoutEnabled{}
	for rows.Next() {
		var flag domain.FlagWithoutEnabled
		if err := rows.Scan(&flag.ID, &flag.Name, &flag.ProjectID); err != nil {
			return nil, err
		}
		flags = append(flags, flag)
	}

	return flags, nil
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

func (r *PgFlagRepository) DeleteFlag(flagId int) error {
	_, err := r.db.Exec(`DELETE FROM flag WHERE id = $1`, flagId)
	return err
}

func (r *PgFlagRepository) GetEnvironmentActiveFlagsByPublicKey(key string) ([]string, error) {
	activeFlagsRows, err := r.db.Query(`
		select f."name" from environment e 
		left join flag_environment fe on e.id = fe.environment
		left join flag f on fe.flag = f.id
		where e.public_key = $1 and fe.enabled = true;
	`, key)

	if err != nil {
		return nil, err
	}

	var activeFlags []string = []string{}
	for activeFlagsRows.Next() {
		var flagName string
		if err := activeFlagsRows.Scan(&flagName); err != nil {
			return nil, err
		}
		activeFlags = append(activeFlags, flagName)
	}

	return activeFlags, nil
}
