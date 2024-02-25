package severs

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/DrumPatiphon/go-rest-api-service/config"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type Isever interface {
	Start()
}

type sever struct {
	app *fiber.App
	cfg config.Iconfig
	db  *sqlx.DB
}

func NewSever(cfg config.Iconfig, db *sqlx.DB) Isever {
	return &sever{
		cfg: cfg,
		db:  db,
		app: fiber.New(fiber.Config{
			AppName:      cfg.App().Name(),
			BodyLimit:    cfg.App().BodyLimit(),
			ReadTimeout:  cfg.App().ReadTimeOut(),
			WriteTimeout: cfg.App().WriteTimeOut(),
			JSONEncoder:  json.Marshal,   // make gofiber faster from goFiber Doc
			JSONDecoder:  json.Unmarshal, // make gofiber faster from goFiber Doc
		}),
	}
}

func (sever *sever) Start() {
	// Middlewares
	middlewares := InitMiddlewares(sever)
	sever.app.Use(middlewares.Logger())
	sever.app.Use(middlewares.Cors())

	// Modlues
	// http://localhost:3000/v1
	v1 := sever.app.Group("v1")

	modules := InitModule(v1, sever, middlewares)

	modules.MonitorModule()
	modules.UserModule()
	modules.AppInfoModule()
	modules.FilesModule()
	modules.ProductsModule()

	sever.app.Use(middlewares.RouterCheck())

	// Graceful Shutdown for safe sever resource and shutdown when have Interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		log.Println("sever is shutting down...")
		_ = sever.app.Shutdown()
	}()

	// Listen to host:port
	log.Println("sever is starting on %v", sever.cfg.App().Url())
	sever.app.Listen(sever.cfg.App().Url())
}
