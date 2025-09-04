package http

import (
	"context"
	"fmt"
	"sms/config"
	"sms/internal/app"

	"sms/docs"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Run(appContainer app.App, cfg config.Server) error {
	router := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
	})

	registerSMSRoutes(appContainer, router)
	docs.SwaggerInfo.Host = ""
	docs.SwaggerInfo.Schemes = []string{}
	docs.SwaggerInfo.BasePath = "/api/v1"

	router.Get("/swagger/*", adaptor.HTTPHandler(httpSwagger.Handler()))

	return router.Listen(fmt.Sprintf(":%d", cfg.Port))
}

func registerSMSRoutes(appContainer app.App, router fiber.Router) {
	ctx := context.Background()
	smsUseCase := appContainer.SMSService(ctx)

	smsHandler := NewSMSHandler(smsUseCase)

	v1 := router.Group("/api/v1")

	// SMS routes
	sms := v1.Group("/sms")
	sms.Post("/", setTraceID(), smsHandler.SendSMS)
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error":   "internal_error",
		"message": err.Error(),
		"code":    code,
	})
}
