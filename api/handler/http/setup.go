package http

import (
	"context"
	"fmt"
	"sms-dispatcher/api/service"
	"sms-dispatcher/app"
	"sms-dispatcher/config"

	"sms-dispatcher/docs"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Run(appContainer *app.App, cfg config.ServerConfig) error {
	router := fiber.New()

	api := router.Group("/api/v1")

	registerSMSAPI(appContainer, cfg, api)

	docs.SwaggerInfo.Host = ""
	docs.SwaggerInfo.Schemes = []string{}
	docs.SwaggerInfo.BasePath = "/api/v1"

	router.Get("/swagger/*", adaptor.HTTPHandler(httpSwagger.Handler()))

	return router.Listen(fmt.Sprintf(":%d", cfg.HttpPort))
}

func registerSMSAPI(appContainer *app.App, cfg config.ServerConfig, router fiber.Router) {
	internalSvc := appContainer.SMSService(context.Background())
	apiSvc := service.NewSMSService(internalSvc)
	router.Post("/sms/send", SendSMSMessage(apiSvc))
	router.Get("/sms/:id", GetSMSMessage(apiSvc))
}
