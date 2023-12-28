package severs

import (
	middlewareHandlers "github.com/DrumPatiphon/go-rest-api-service/modules/middleware/middlewareHandlers"
	middlewareRepositories "github.com/DrumPatiphon/go-rest-api-service/modules/middleware/middlewareRepositories"
	middlewareUsecases "github.com/DrumPatiphon/go-rest-api-service/modules/middleware/middlewareUsecases"
	mornitorHandlers "github.com/DrumPatiphon/go-rest-api-service/modules/monitor/monitorHandlers"
	"github.com/DrumPatiphon/go-rest-api-service/modules/users/usersHandlers"
	"github.com/DrumPatiphon/go-rest-api-service/modules/users/usersRepositories"
	"github.com/DrumPatiphon/go-rest-api-service/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
	UserModule()
}

type moduleFactory struct {
	router     fiber.Router
	sever      *sever
	middleware middlewareHandlers.ImiddlewareHandler
}

func InitModule(router fiber.Router, sever *sever, middleware middlewareHandlers.ImiddlewareHandler) IModuleFactory {
	return &moduleFactory{
		router:     router,
		sever:      sever,
		middleware: middleware,
	}
}

func InitMiddlewares(sever *sever) middlewareHandlers.ImiddlewareHandler {
	repository := middlewareRepositories.Middlewarerepository(sever.db)
	usecase := middlewareUsecases.MiddlewareUsecase(repository)
	handler := middlewareHandlers.MiddlewareHandler(sever.cfg, usecase)
	return handler
}

func (module *moduleFactory) MonitorModule() {
	handler := mornitorHandlers.MonitorHandler(module.sever.cfg)

	module.router.Get("/", handler.HelthCheck)
}

func (module *moduleFactory) UserModule() {
	repository := usersRepositories.UserRepository(module.sever.db)
	usecase := usersUsecases.UserUsecases(module.sever.cfg, repository)
	handler := usersHandlers.UserHandler(module.sever.cfg, usecase)

	// /v1/users/sign
	router := module.router.Group("/users")

	router.Post("/signup", handler.SignUpCustomer)
	router.Post("/signIn", handler.SignIn)
	router.Post("/refresh", handler.RefreshPassport)
	router.Post("/signout", handler.SignOut)
	router.Post("/signup-admin", handler.SignOut)

	router.Get("/:user_id", module.middleware.JwtAuth(), module.middleware.ParamsCheck(), handler.GetUserProfile)
	router.Get("/admin/secret", module.middleware.JwtAuth(), module.middleware.Autorize(2), handler.GenerateAdminToken)

	// Initial admin ขึ้นมา 1 คนใน Database (insert ใน sql)
	// Generate Admin Key
	// ทุกครั้งที่ทำการสมัคร admin เพิ่ม ให้ส่ง Admin Token มาด้วยทุกครั้ง ผ่าน Middleware
}
