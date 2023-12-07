package usersHandlers

import (
	"github.com/DrumPatiphon/go-rest-api-service/config"
	"github.com/DrumPatiphon/go-rest-api-service/modules/entities"
	"github.com/DrumPatiphon/go-rest-api-service/modules/users"
	"github.com/DrumPatiphon/go-rest-api-service/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type userHandlerErrCode string

const (
	signUpCustomerErr userHandlerErrCode = "users-001"
)

type IUserUsecases interface {
	SignUpCustomer(c *fiber.Ctx) error
}

type usersHandler struct {
	cfg          config.Iconfig
	usersUsecase usersUsecases.IUserUsecases
}

func UserHandler(cfg config.Iconfig, usersUsecase usersUsecases.IUserUsecases) IUserUsecases {
	return &usersHandler{
		cfg:          cfg,
		usersUsecase: usersUsecase,
	}
}

func (h *usersHandler) SignUpCustomer(c *fiber.Ctx) error {
	// Req body parser
	req := new(users.UserRegisterReq)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code, //400
			string(signUpCustomerErr),
			err.Error(),
		).Res()
	}

	// Email validation
	if !req.IsEmail() {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code, //400
			string(signUpCustomerErr),
			"email pattern is invalid",
		).Res()
	}

	// Insert
	result, err := h.usersUsecase.InsertCustomer(req)
	if err != nil {
		switch err.Error() {
		case "username has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code, //400
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		case "email has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code, //400
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code, //500
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		}
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}
