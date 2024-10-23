package infrastructure

import (
	"database/sql"

	entity "github.com/Palma99/feature-flag-service/internals/domain/entity"
)

type PgUserRepository struct {
	db *sql.DB
}

func NewPgUserRepository(db *sql.DB) *PgUserRepository {
	return &PgUserRepository{
		db,
	}
}

func (r *PgUserRepository) CreateUser(user *entity.User) error {
	_, err := r.db.Exec(`
		INSERT INTO users (nickname, email, "password") VALUES ($1, $2, $3)
	`, user.Nickname, user.Email, user.Password)

	if err != nil {
		return err
	}

	return nil
}

func (r *PgUserRepository) GetUserByEmail(email string) (*entity.User, error) {
	userRow := r.db.QueryRow(`
		SELECT 
			id,
			nickname,
			email,
			"password"
		FROM users 
		WHERE email = $1
	`, email)

	if userRow.Err() != nil {
		return nil, userRow.Err()
	}

	user := entity.User{}

	if err := userRow.Scan(
		&user.ID,
		&user.Nickname,
		&user.Email,
		&user.Password,
	); err != nil {
		return nil, err
	}

	return &user, nil
}
