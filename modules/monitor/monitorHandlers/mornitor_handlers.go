package mornitorHandlers

import (
	"github.com/DrumPatiphon/go-rest-api-service/config"
	"github.com/DrumPatiphon/go-rest-api-service/modules/entities"
	"github.com/DrumPatiphon/go-rest-api-service/modules/monitor"
	"github.com/gofiber/fiber/v2"
)

// มีหน้าที่รับ api req จาก network หรือ protocol

type IMornitorHandler interface {
	HelthCheck(context *fiber.Ctx) error //handler รับ param เป็น fiber.context เท่านั้น
}

type monitorHandler struct {
	config config.Iconfig
}

func MonitorHandler(config config.Iconfig) IMornitorHandler {
	return &monitorHandler{
		config: config,
	}
}

func (handler *monitorHandler) HelthCheck(context *fiber.Ctx) error {
	res := &monitor.Mornitor{
		Name:    handler.config.App().Name(),
		Version: handler.config.App().Version(),
	}
	return entities.NewResponse(context).Success(fiber.StatusOK, res).Res()
}
