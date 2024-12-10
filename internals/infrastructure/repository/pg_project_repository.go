package infrastructure

import (
	"database/sql"
	"encoding/json"

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

func (r *PgProjectRepository) GetUserProjectRelation(userId int, projectId int64) (*int64, error) {
	row := r.db.QueryRow(`
	select id from users_project up 
		where up.user_id = $1 and up.project_id = $2
	`, userId, projectId)

	var id int64
	if err := row.Scan(&id); err != nil {
		return nil, err
	}

	return &id, nil
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

func (r *PgProjectRepository) GetProjectDetails(projectId int64) (*entity.ProjectWithMembers, error) {
	row := r.db.QueryRow(`
		SELECT
			p.id AS "projectId",
			p.owner_id as "ownerId",
			p.name as "name",
			COALESCE(
				jsonb_agg(DISTINCT jsonb_build_object(
					'id', up.user_id,
					'nickname', u.nickname,
					'email', u.email
				)) FILTER (WHERE up.user_id IS NOT NULL),
				'[]'::jsonb
			) AS "members",
			COALESCE(
				jsonb_agg(DISTINCT jsonb_build_object(
					'id', e.id,
					'name', e.name,
					'public_key', e.public_key,
					'project_id', e.project_id
				)) FILTER (WHERE e.id IS NOT NULL),
				'[]'::jsonb
			) AS "environments"
		FROM project AS p
		LEFT JOIN users_project up ON up.project_id = p.id
		LEFT JOIN users u ON u.id = up.user_id
		LEFT JOIN environment e ON e.project_id = p.id
		WHERE p.id = $1
		GROUP BY p.id;
	`, projectId)

	var projectID int
	var ownerID int
	var projectName string
	var membersJSON []byte
	var environmentsJSON []byte

	err := row.Scan(&projectID, &ownerID, &projectName, &membersJSON, &environmentsJSON)
	if err != nil {
		return nil, err
	}

	var members []entity.User
	if err := json.Unmarshal(membersJSON, &members); err != nil {
		return nil, err
	}

	var environments []entity.Environment
	if err := json.Unmarshal(environmentsJSON, &environments); err != nil {
		return nil, err
	}

	project := entity.ProjectWithMembers{
		Project: entity.Project{
			ID:           int64(projectID),
			OwnerId:      ownerID,
			Name:         projectName,
			Environments: environments,
		},
		Members: members,
	}

	return &project, nil
}
