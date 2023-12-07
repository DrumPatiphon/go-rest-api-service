package usersRepositories

import (
	"github.com/DrumPatiphon/go-rest-api-service/modules/users"
	"github.com/DrumPatiphon/go-rest-api-service/modules/users/usersPatterns"
	"github.com/jmoiron/sqlx"
)

type IUserRepository interface {
	InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error)
}

type usersRepository struct {
	db *sqlx.DB
}

func UserRepositories(db *sqlx.DB) IUserRepository {
	return &usersRepository{
		db: db,
	}
}

func (r *usersRepository) InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error) {
	result := usersPatterns.InsertUser(r.db, req, isAdmin)

	var err error

	if isAdmin {
		result, err = result.Admin()
		if err != nil {
			return nil, err
		}
	} else {
		result, err = result.Customer()
		if err != nil {
			return nil, err
		}
	}

	// Get Result From inserting
	user, err := result.Result()
	if err != nil {
		return nil, err
	}
	return user, nil
}
