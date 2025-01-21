package infra

import (
	"database/sql"
	"github.com/killinsun/codetest-docker/app/domain"
)

type RdbUserRepository struct {
	conn *sql.DB
}

func NewUserRepository(conn *sql.DB) *RdbUserRepository {
	return &RdbUserRepository{conn: conn}
}

func (r *RdbUserRepository) FindByApiKey(apiKey string) (*domain.User, error) {
	row := r.conn.QueryRow("SELECT id, name, api_key FROM users WHERE api_key = ?", apiKey)
	user := &domain.User{}

	if err := row.Scan(&user.ID, &user.Name, &user.ApiKey); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}
