package middlewarehandlers

import (
	"strings"

	"github.com/DrumPatiphon/go-rest-api-service/config"
	"github.com/DrumPatiphon/go-rest-api-service/modules/entities"
	middlewareUsecases "github.com/DrumPatiphon/go-rest-api-service/modules/middleware/middlewareUsecases"
	"github.com/DrumPatiphon/go-rest-api-service/pkg/serviceauth"
	"github.com/DrumPatiphon/go-rest-api-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type middlewareHandlersErrCode string

const (
	routerCheckErr middlewareHandlersErrCode = "middleware-001"
	jwtAuthErr     middlewareHandlersErrCode = "middleware-002"
	paramCheckErr  middlewareHandlersErrCode = "middleware-003"
	authorizeErr   middlewareHandlersErrCode = "middleware-004"
	apiKeyErr      middlewareHandlersErrCode = "middleware-005"
)

type ImiddlewareHandler interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
	JwtAuth() fiber.Handler
	ParamsCheck() fiber.Handler
	Autorize(expectRoleId ...int) fiber.Handler
	ApiKeyAuth() fiber.Handler
}

type middlewaresHandler struct {
	cfg               config.Iconfig
	middlewareUsecase middlewareUsecases.ImiddlewareUsecase
}

func MiddlewareHandler(cfg config.Iconfig, middlewareUsecase middlewareUsecases.ImiddlewareUsecase) ImiddlewareHandler {
	return &middlewaresHandler{
		cfg:               cfg,
		middlewareUsecase: middlewareUsecase,
	}
}

func (handler *middlewaresHandler) Cors() fiber.Handler {
	return cors.New(cors.Config{
		Next:             cors.ConfigDefault.Next,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	})
}

func (handler *middlewaresHandler) RouterCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return entities.NewResponse(c).Error(
			fiber.ErrNotFound.Code,
			string(routerCheckErr),
			"router not found",
		).Res()
	}
}

func (handler *middlewaresHandler) Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} [${ip}] ${status} - ${method} ${path}\n",
		TimeFormat: "02/01/2006",
		TimeZone:   "Bankok/Asia",
	})
}

func (h *middlewaresHandler) JwtAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		result, err := serviceauth.ParseToken(h.cfg.Jwt(), token)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErr),
				err.Error(),
			).Res()
		}

		claims := result.Claims
		if !h.middlewareUsecase.FindAccessToken(claims.Id, token) {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErr),
				"no permission to access",
			).Res()
		}

		// Set UserId
		c.Locals("userId", claims.Id)
		c.Locals("userRoleId", claims.RoleId)
		return c.Next()
	}
}

func (h *middlewaresHandler) ParamsCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("userId")
		if c.Params("user_id") != userId {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(paramCheckErr),
				"never gonna give you up",
			).Res()
		}
		return c.Next()
	}
}

func (h *middlewaresHandler) Autorize(expectRoleId ...int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoleId, ok := c.Locals("userRoleId").(int)
		if !ok {
			entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(authorizeErr),
				"role_id is not int type",
			).Res()
		}

		roles, err := h.middlewareUsecase.FindRole()
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(authorizeErr),
				err.Error(),
			).Res()
		}

		sum := 0
		for _, v := range expectRoleId {
			sum += v
		}

		expectValueBinary := utils.BinaryConverter(sum, len(roles))
		userValueBinary := utils.BinaryConverter(userRoleId, len(roles))

		for i := range userValueBinary {
			if userValueBinary[i]&expectValueBinary[i] == 1 {
				return c.Next()
			}
		}

		return entities.NewResponse(c).Error(
			fiber.ErrUnauthorized.Code,
			string(authorizeErr),
			"no permission to access",
		).Res()
	}
}

func (h *middlewaresHandler) ApiKeyAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := c.Get("x-Api-key")
		if _, err := serviceauth.ParseApiKey(h.cfg.Jwt(), key); err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(apiKeyErr),
				"api Key is invalid or requried",
			).Res()
		}
		return c.Next()
	}
}
