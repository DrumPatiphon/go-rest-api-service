package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/DrumPatiphon/go-rest-api-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type ILogger interface {
	Print() ILogger
	Save()
	SetQuery(context *fiber.Ctx)
	SetBody(context *fiber.Ctx)
	SetResponse(res any)
	// http:localhost:300/v1/products
}

type logger struct {
	Time       string `json:"time"`
	Ip         string `json:"ip"`
	Method     string `json:"method"`
	StatusCode int    `json:"status_code"`
	Path       string `json:"path"`
	Query      any    `json:"query"`
	Body       any    `json:"body"`
	Response   any    `json:"response"`
}

func InitLogger(context *fiber.Ctx, res any) ILogger {
	log := &logger{
		Time:       time.Now().Local().Format("2006-01-02 15:04:05"),
		Ip:         context.IP(),
		Method:     context.Method(),
		Path:       context.Path(),
		StatusCode: context.Response().StatusCode(),
	}
	log.SetQuery(context)
	log.SetBody(context)
	log.SetResponse(res)
	return log
}

func (l *logger) Print() ILogger {
	utils.Debug(l)
	return l
}

func (l *logger) Save() {
	data := utils.Output(l)

	filename := fmt.Sprintf("./assets/logs/logger_%v.txt", strings.ReplaceAll(time.Now().Format("2006-01-02"), "-", ""))
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) // 0666 is UNIX code of permission meaning is can read erite create but can't delete
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()
	file.WriteString(string(data) + "\n")
}

func (l *logger) SetQuery(context *fiber.Ctx) {
	var body any
	if err := context.QueryParser(&body); err != nil {
		log.Printf("query parser error: %v", err)
	}
	l.Query = body
}

func (l *logger) SetBody(context *fiber.Ctx) {
	var body any
	if err := context.BodyParser(&body); err != nil {
		log.Printf("body parser error: %v", err)
	}

	switch l.Path {
	case "v1/users/signup":
		l.Body = "never gonna give you up"
	default:
		l.Body = body
	}
}

func (l *logger) SetResponse(res any) {
	l.Response = res
}
