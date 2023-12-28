package middlewareusecases

import (
	"github.com/DrumPatiphon/go-rest-api-service/modules/middleware"
	middlewareRepositories "github.com/DrumPatiphon/go-rest-api-service/modules/middleware/middlewareRepositories"
)

type ImiddlewareUsecase interface {
	FindAccessToken(userId, accessToken string) bool
	FindRole() ([]*middleware.Role, error)
}

type middlewaUsecase struct {
	middlewareRepository middlewareRepositories.ImiddlewareRepository
}

func MiddlewareUsecase(middlewareRepository middlewareRepositories.ImiddlewareRepository) ImiddlewareUsecase {
	return &middlewaUsecase{
		middlewareRepository: middlewareRepository,
	}
}

func (u *middlewaUsecase) FindAccessToken(userId, accessToken string) bool {
	return u.middlewareRepository.FindAccessToken(userId, accessToken)
}

func (u *middlewaUsecase) FindRole() ([]*middleware.Role, error) {
	roles, err := u.middlewareRepository.FindRole()
	if err != nil {
		return nil, err
	}
	return roles, nil
}
