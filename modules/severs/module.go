package severs

import (
	appinfoHandlers "github.com/DrumPatiphon/go-rest-api-service/modules/appInfo/appInfoHandlers"
	appinfoRepositories "github.com/DrumPatiphon/go-rest-api-service/modules/appInfo/appInfoRepositories"
	appinfoUsecases "github.com/DrumPatiphon/go-rest-api-service/modules/appInfo/appInfoUsecases"
	"github.com/DrumPatiphon/go-rest-api-service/modules/files/filesHandlers"
	"github.com/DrumPatiphon/go-rest-api-service/modules/files/filesUsecases"
	middlewareHandlers "github.com/DrumPatiphon/go-rest-api-service/modules/middleware/middlewareHandlers"
	middlewareRepositories "github.com/DrumPatiphon/go-rest-api-service/modules/middleware/middlewareRepositories"
	middlewareUsecases "github.com/DrumPatiphon/go-rest-api-service/modules/middleware/middlewareUsecases"
	mornitorHandlers "github.com/DrumPatiphon/go-rest-api-service/modules/monitor/monitorHandlers"
	"github.com/DrumPatiphon/go-rest-api-service/modules/products/productsHandlers"
	"github.com/DrumPatiphon/go-rest-api-service/modules/products/productsRepositories"
	"github.com/DrumPatiphon/go-rest-api-service/modules/products/productsUsecases"
	"github.com/DrumPatiphon/go-rest-api-service/modules/users/usersHandlers"
	"github.com/DrumPatiphon/go-rest-api-service/modules/users/usersRepositories"
	"github.com/DrumPatiphon/go-rest-api-service/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
	UserModule()
	AppInfoModule()
	FilesModule()
	ProductsModule()
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

	router.Post("/signup", module.middleware.ApiKeyAuth(), handler.SignUpCustomer)
	router.Post("/signIn", module.middleware.ApiKeyAuth(), handler.SignIn)
	router.Post("/refresh", module.middleware.ApiKeyAuth(), handler.RefreshPassport)
	router.Post("/signout", module.middleware.ApiKeyAuth(), handler.SignOut)
	router.Post("/signup-admin", module.middleware.JwtAuth(), module.middleware.Autorize(2), handler.SingUpAdmin)

	router.Get("/:user_id", module.middleware.JwtAuth(), module.middleware.ParamsCheck(), handler.GetUserProfile)
	router.Get("/admin/secret", module.middleware.JwtAuth(), module.middleware.Autorize(2), handler.GenerateAdminToken)

	// Initial admin ขึ้นมา 1 คนใน Database (insert ใน sql)
	// Generate Admin Key
	// ทุกครั้งที่ทำการสมัคร admin เพิ่ม ให้ส่ง Admin Token มาด้วยทุกครั้ง ผ่าน Middleware
}

func (module *moduleFactory) AppInfoModule() {
	repository := appinfoRepositories.AppInfoRepository(module.sever.db)
	usecase := appinfoUsecases.AppInfoUsecase(repository)
	handler := appinfoHandlers.AppInfoHandler(module.sever.cfg, usecase)

	router := module.router.Group("/appinfo")

	router.Post("/categories", module.middleware.JwtAuth(), module.middleware.Autorize(2), handler.AddCategory)

	router.Get("/categories", module.middleware.ApiKeyAuth(), handler.FindCategory)
	router.Get("/apikey", module.middleware.JwtAuth(), module.middleware.Autorize(2), handler.GenerateApiKey)

	router.Delete("/:category_id/categories", module.middleware.JwtAuth(), module.middleware.Autorize(2), handler.RemoveCategory)
}

func (m *moduleFactory) FilesModule() {
	usecase := filesUsecases.FilesUsecase(m.sever.cfg)
	handler := filesHandlers.FilesHandler(m.sever.cfg, usecase)

	router := m.router.Group("/files")

	router.Post("/upload", m.middleware.JwtAuth(), m.middleware.Autorize(2), handler.UploadFile)
	router.Patch("/delete", m.middleware.JwtAuth(), m.middleware.Autorize(2), handler.DeleteFile) //ใช้ patch จะได้เพิ่ม body
}

func (m *moduleFactory) ProductsModule() {
	fileUsecase := filesUsecases.FilesUsecase(m.sever.cfg)

	productsRepository := productsRepositories.ProductRepository(m.sever.db, m.sever.cfg, fileUsecase)
	productsUsecase := productsUsecases.ProductsUsecases(productsRepository)
	productsHandler := productsHandlers.ProductsHandler(m.sever.cfg, productsUsecase, fileUsecase)

	router := m.router.Group("/products")

	router.Post("/", m.middleware.JwtAuth(), m.middleware.Autorize(2), productsHandler.InsertProduct)
	router.Patch("/:product_id", m.middleware.JwtAuth(), m.middleware.Autorize(2), productsHandler.UpdateProduct)

	router.Get("/", m.middleware.ApiKeyAuth(), productsHandler.FindProduct)
	router.Get("/:product_id", m.middleware.ApiKeyAuth(), productsHandler.FindOneProduct)

	router.Delete("/:product_id", m.middleware.JwtAuth(), m.middleware.Autorize(2), productsHandler.DeleteProduct)
}
