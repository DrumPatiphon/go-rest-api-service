package middlewareusecases

import middlewareRepositories "github.com/DrumPatiphon/go-rest-api-service/modules/users/middleware/middlewareRepositories"

type ImiddlewareUsecase interface {
}

type middlewaUsecase struct {
	middlewareRepository middlewareRepositories.ImiddlewareRepository
}

func MiddlewareUsecase(middlewareRepository middlewareRepositories.ImiddlewareRepository) ImiddlewareUsecase {
	return &middlewaUsecase{
		middlewareRepository: middlewareRepository,
	}
}
