package productsHandlers

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/DrumPatiphon/go-rest-api-service/config"
	"github.com/DrumPatiphon/go-rest-api-service/modules/appInfo"
	"github.com/DrumPatiphon/go-rest-api-service/modules/entities"
	"github.com/DrumPatiphon/go-rest-api-service/modules/files/filesUsecases"
	"github.com/DrumPatiphon/go-rest-api-service/modules/products"
	"github.com/DrumPatiphon/go-rest-api-service/modules/products/productsUsecases"
	"github.com/gofiber/fiber/v2"
)

type productsHandlerErrCode string

const (
	findOneProductErr productsHandlerErrCode = "products-001"
	findProductErr    productsHandlerErrCode = "products-002"
	insertProductErr  productsHandlerErrCode = "products-003"
	updateProductErr  productsHandlerErrCode = "products-004"
)

type IProductHandler interface {
	FindOneProduct(c *fiber.Ctx) error
	FindProduct(c *fiber.Ctx) error
	InsertProduct(c *fiber.Ctx) error
	UpdateProduct(c *fiber.Ctx) error
}

type productsHandler struct {
	cfg              config.Iconfig
	productsUsecases productsUsecases.IProductUseCase
	filesUsecases    filesUsecases.IFilesUsecases
}

func ProductsHandler(cfg config.Iconfig, productsUsecases productsUsecases.IProductUseCase, filesUsecases filesUsecases.IFilesUsecases) IProductHandler {
	return &productsHandler{
		cfg:              cfg,
		productsUsecases: productsUsecases,
		filesUsecases:    filesUsecases,
	}
}

func (h *productsHandler) FindOneProduct(c *fiber.Ctx) error {
	productId := strings.Trim(c.Params("product_id"), " ")

	product, err := h.productsUsecases.FindOneProduct(productId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findOneProductErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, product).Res()
}

func (h *productsHandler) FindProduct(c *fiber.Ctx) error {
	req := &products.ProductFilter{
		PaginationReq: &entities.PaginationReq{},
		SortReq:       &entities.SortReq{},
	}

	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findProductErr),
			err.Error(),
		).Res()
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 5 {
		req.Limit = 5
	}

	if req.OrderBy == "" {
		req.OrderBy = "title"
	}
	if req.Sort == "" {
		req.Sort = "ASC"
	}

	products := h.productsUsecases.FindProduct(req)
	return entities.NewResponse(c).Success(fiber.StatusOK, products).Res()
}

func (h *productsHandler) InsertProduct(c *fiber.Ctx) error {
	req := &products.Product{
		Category: &appInfo.Category{},
		Images:   make([]*entities.Image, 0),
	}
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertProductErr),
			err.Error(),
		).Res()
	}

	if req.Category.Id <= 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertProductErr),
			"cagetory is is invalid",
		).Res()
	}

	product, err := h.productsUsecases.InsertProduct(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(insertProductErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, product).Res()
}

func (h *productsHandler) UpdateProduct(c *fiber.Ctx) error {
	productId := strings.Trim(c.Params("product_id"), " ")

	req := &products.Product{
		Images:   make([]*entities.Image, 0),
		Category: &appInfo.Category{},
	}
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateProductErr),
			err.Error(),
		).Res()
	}
	req.Id = productId

	product, err := h.productsUsecases.UpdateProduct(req)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		errMsg := fmt.Sprintf("%s:%d %s", file, line, err.Error())
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(updateProductErr),
			errMsg,
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, product).Res()
}
