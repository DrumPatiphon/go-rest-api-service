package appinfoHandlers

import (
	"strconv"
	"strings"

	"github.com/DrumPatiphon/go-rest-api-service/config"
	"github.com/DrumPatiphon/go-rest-api-service/modules/appInfo"
	appinfoUsecases "github.com/DrumPatiphon/go-rest-api-service/modules/appInfo/appInfoUsecases"
	"github.com/DrumPatiphon/go-rest-api-service/modules/entities"
	"github.com/DrumPatiphon/go-rest-api-service/pkg/serviceauth"
	"github.com/gofiber/fiber/v2"
)

type appInfoHandlerErrcode string

const (
	generrateApiKeyErr appInfoHandlerErrcode = "appinfo-001"
	findCategoryErr    appInfoHandlerErrcode = "appinfo-002"
	addCategoryErr     appInfoHandlerErrcode = "appinfo-003"
	RemoveCategoryErr  appInfoHandlerErrcode = "appinfo-004"
)

type IAppInfoHandler interface {
	GenerateApiKey(c *fiber.Ctx) error
	FindCategory(c *fiber.Ctx) error
	AddCategory(c *fiber.Ctx) error
	RemoveCategory(c *fiber.Ctx) error
}

type appinfoHandler struct {
	cfg            config.Iconfig
	appInfousecase appinfoUsecases.IAppInfoUsecase
}

func AppInfoHandler(cfg config.Iconfig, appInfousecase appinfoUsecases.IAppInfoUsecase) IAppInfoHandler {
	return &appinfoHandler{
		cfg:            cfg,
		appInfousecase: appInfousecase,
	}
}

func (h *appinfoHandler) GenerateApiKey(c *fiber.Ctx) error {
	apiKey, err := serviceauth.NewServiceAuth(
		serviceauth.ApiKey,
		h.cfg.Jwt(),
		nil,
	)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(generrateApiKeyErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		&struct {
			Key string `json:"key"`
		}{
			Key: apiKey.SignToken(),
		},
	).Res()
}

func (h *appinfoHandler) FindCategory(c *fiber.Ctx) error {
	req := new(appInfo.CategoryFilter)
	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findCategoryErr),
			err.Error(),
		).Res()
	}

	category, err := h.appInfousecase.FindCategory(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findCategoryErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		category,
	).Res()
}

func (h *appinfoHandler) AddCategory(c *fiber.Ctx) error {
	req := make([]*appInfo.Category, 0)
	if err := c.BodyParser(&req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(addCategoryErr),
			err.Error(),
		).Res()
	}
	if len(req) == 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(addCategoryErr),
			"categories request are emty",
		).Res()
	}

	if err := h.appInfousecase.InsertCagetory(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(addCategoryErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, req).Res()
}

func (h *appinfoHandler) RemoveCategory(c *fiber.Ctx) error {

	categoryId := strings.Trim(c.Params("category_id"), " ")
	categoryIdInt, err := strconv.Atoi(categoryId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(RemoveCategoryErr),
			"id type is invalid",
		).Res()
	}
	if categoryIdInt <= 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(RemoveCategoryErr),
			"id is must more than zero",
		).Res()
	}

	if err := h.appInfousecase.DeleteCategory(categoryIdInt); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(RemoveCategoryErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		&struct {
			CategoryId int `json:"category_id"`
		}{
			CategoryId: categoryIdInt,
		},
	).Res()
}
