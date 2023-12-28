package middlewarerepositories

import (
	"fmt"

	"github.com/DrumPatiphon/go-rest-api-service/modules/middleware"
	"github.com/jmoiron/sqlx"
)

type ImiddlewareRepository interface {
	FindAccessToken(userId, accessToken string) bool
	FindRole() ([]*middleware.Role, error)
}

type middlewareRepository struct {
	db *sqlx.DB
}

func Middlewarerepository(db *sqlx.DB) ImiddlewareRepository {
	return &middlewareRepository{
		db: db,
	}
}

func (r *middlewareRepository) FindAccessToken(userId string, accessToken string) bool {
	query := `
	SELECT
		CASE WHEN COUNT(id) = 1 THEN TRUE ELSE FALSE END
	FROM oauth
	WHERE user_id = $1
	AND access_token = $2`

	var check bool
	if err := r.db.Get(&check, query, userId, accessToken); err != nil {
		return false
	}
	return true
}

func (r *middlewareRepository) FindRole() ([]*middleware.Role, error) {
	query := `
	SELECT
			id,
			title
	FROM roles
	ORDER BY id DESC;`

	roles := make([]*middleware.Role, 0)
	if err := r.db.Select(&roles, query); err != nil {
		return nil, fmt.Errorf("roles are empty : %v", err.Error())
	}
	return roles, nil
}
