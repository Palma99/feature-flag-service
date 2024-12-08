package infrastructure

import (
	"database/sql"

	entity "github.com/Palma99/feature-flag-service/internals/domain/entity"
	domain "github.com/Palma99/feature-flag-service/internals/domain/repository"
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
		SELECT id, name, public_key FROM environment WHERE project_id = $1
	`, id)

	if err != nil {
		return nil, err
	}

	for environmentRows.Next() {
		environment := &entity.Environment{}
		if err := environmentRows.Scan(
			environment.ID, environment.Name, environment.PublicKey,
		); err != nil {
			return nil, err
		}
		project.Environments = append(project.Environments, *environment)
	}

	return project, nil
}

func (r *PgProjectRepository) CreateProject(project *domain.CreateProjectDTO) (*entity.Project, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	createProjectQuery := `
		INSERT INTO project (name, owner_id) VALUES ($1, $2) RETURNING id
	`
	stmt, err := tx.Prepare(createProjectQuery)
	if err != nil {
		return nil, err
	}

	var createdProjectId int64

	err = stmt.QueryRow(project.Name, project.OwnerId).Scan(&createdProjectId)
	if err != nil {
		return nil, err
	}

	createdProject := &entity.Project{
		ID:      createdProjectId,
		Name:    project.Name,
		OwnerId: project.OwnerId,
	}

	_, err = tx.Exec(`
		INSERT INTO users_project (user_id, project_id) VALUES ($1, $2)
	`, project.OwnerId, createdProjectId)

	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return createdProject, nil
}

func (r *PgProjectRepository) GetProjectsByUserId(userId int) ([]entity.Project, error) {
	projectRows, err := r.db.Query(`
		SELECT p.id, p.name, p.owner_id FROM project AS p
			LEFT JOIN users_project up ON up.project_id = p.id
			WHERE up.user_id = $1
		`, userId)

	if err != nil {
		return nil, err
	}

	projects := []entity.Project{}

	for projectRows.Next() {
		project := entity.Project{}
		if err := projectRows.Scan(&project.ID, &project.Name, &project.OwnerId); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, nil
}

func (r *PgProjectRepository) GetUserProjectByProjectId(userId int, projectId int64) (*entity.Project, error) {

	rows, err := r.db.Query(`
		select p.id as id, p.name as name, owner_id, e.name as env_nam, e.id as env_id FROM project as p
			left join users_project up on up.project_id = p.id
			left join environment e on e.project_id = p.id
			where up.user_id = $1 and p.id = $2
		`, userId, projectId)

	if err != nil {
		return nil, err
	}

	project := &entity.Project{}

	for rows.Next() {
		environment := &entity.Environment{
			ProjectID: projectId,
		}
		if err := rows.Scan(&project.ID, &project.Name, &project.OwnerId, &environment.Name, &environment.ID); err != nil {
			return nil, err
		}
		project.Environments = append(project.Environments, *environment)
	}

	return project, nil
}
