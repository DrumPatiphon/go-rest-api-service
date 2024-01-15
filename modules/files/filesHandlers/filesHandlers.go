package filesHandlers

import (
	"fmt"
	"math"
	"path/filepath"
	"strings"

	"github.com/DrumPatiphon/go-rest-api-service/config"
	"github.com/DrumPatiphon/go-rest-api-service/modules/entities"
	"github.com/DrumPatiphon/go-rest-api-service/modules/files"
	"github.com/DrumPatiphon/go-rest-api-service/modules/files/filesUsecases"
	"github.com/DrumPatiphon/go-rest-api-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type filesHandlersErrCode string

const (
	uploadErr      filesHandlersErrCode = "files-001"
	deletefilesErr filesHandlersErrCode = "files-002"
)

type IFilesHandler interface {
	UploadFile(c *fiber.Ctx) error
	DeleteFile(c *fiber.Ctx) error
}

type filesHandler struct {
	cfg         config.Iconfig
	fileUsecase filesUsecases.IFilesUsecases
}

func FilesHandler(cfg config.Iconfig, fileUsecase filesUsecases.IFilesUsecases) IFilesHandler {
	return &filesHandler{
		cfg:         cfg,
		fileUsecase: fileUsecase,
	}
}

func (h *filesHandler) UploadFile(c *fiber.Ctx) error {
	req := make([]*files.FileReq, 0)

	form, err := c.MultipartForm()
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(uploadErr),
			err.Error(),
		).Res()
	}

	filesReq := form.File["files"]
	destiantion := c.FormValue("destination")

	// Files ext validation
	extMap := map[string]string{
		"png":  "png",
		"jpg":  "jpg",
		"jpeg": "jpeg",
	}
	for _, file := range filesReq {
		ext := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
		if extMap[ext] != ext || extMap[ext] == "" {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(uploadErr),
				"extension is not acceptable",
			).Res()
		}

		if file.Size > int64(h.cfg.App().FileLimit()) {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(uploadErr),
				fmt.Sprint("file size must less than %d MB", int(math.Ceil(float64(h.cfg.App().FileLimit())/math.Pow(1024, 2)))),
			).Res()
		}

		filename := utils.RandomFilename(ext)
		req = append(req, &files.FileReq{
			File:        file,
			Destination: destiantion + "/" + filename,
			FileName:    filename,
			Extension:   ext,
		})
	}

	res, err := h.fileUsecase.UploadToGCP(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(uploadErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, res).Res()
}

func (h *filesHandler) DeleteFile(c *fiber.Ctx) error {
	req := make([]*files.DeleteFileReq, 0)
	if err := c.BodyParser(&req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(deletefilesErr),
			err.Error(),
		).Res()
	}

	if err := h.fileUsecase.DeleteFileOnGCP(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deletefilesErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, nil).Res()
}
