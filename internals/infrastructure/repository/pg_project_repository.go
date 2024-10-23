package infrastructure

import (
	"database/sql"

	entity "github.com/Palma99/feature-flag-service/internals/domain/entity"
)

type PgProjectRepository struct {
	db *sql.DB
}

func NewPgProjectRepository(db *sql.DB) *PgProjectRepository {
	return &PgProjectRepository{
		db,
	}
}

func (r *PgProjectRepository) GetProjectById(id int) (*entity.Project, error) {
	// get the project by id

	projectRow := r.db.QueryRow(`
		SELECT id, name FROM project WHERE id = $1
	`, id)

	if projectRow.Err() != nil {
		return nil, projectRow.Err()
	}

	project := &entity.Project{}

	if err := projectRow.Scan(project.ID, project.Name); err != nil {
		return nil, err
	}

	environmentRows, err := r.db.Query(`
		SELECT id, name, public_key, private_key FROM environment WHERE project_id = $1
	`, id)

	if err != nil {
		return nil, err
	}

	for environmentRows.Next() {
		environment := &entity.Environment{}
		if err := environmentRows.Scan(
			environment.ID, environment.Name, environment.PublicKey, environment.PrivateKey,
		); err != nil {
			return nil, err
		}
		project.Environments = append(project.Environments, *environment)
	}

	return project, nil
}

func (r *PgProjectRepository) CreateProject(project *entity.Project) error {
	_, err := r.db.Exec(`
		INSERT INTO project (name) VALUES ($1)	
	`, project.Name)

	if err != nil {
		return err
	}

	return nil
}
