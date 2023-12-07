package usersUsecases

import (
	"github.com/DrumPatiphon/go-rest-api-service/config"
	"github.com/DrumPatiphon/go-rest-api-service/modules/users"
	"github.com/DrumPatiphon/go-rest-api-service/modules/users/usersRepositories"
)

type IUserUsecases interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
}

type usersUsecases struct {
	cfg               config.Iconfig
	usersRepositories usersRepositories.IUserRepository
}

func UserUsecases(cfg config.Iconfig, usersRepositories usersRepositories.IUserRepository) IUserUsecases {
	return &usersUsecases{
		cfg:               cfg,
		usersRepositories: usersRepositories,
	}
}

func (u *usersUsecases) InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error) {
	// Hashing a password
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	// Insert user
	result, err := u.usersRepositories.InsertUser(req, false)
	if err != nil {
		return nil, err
	}
	return result, nil
}
