package http

import (
	"context"
	"fmt"
	"sms/config"
	"sms/internal/app"

	"github.com/gofiber/fiber/v2"
)

func Run(appContainer app.App, cfg config.Server) error {
	router := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
	})

	registerSMSRoutes(appContainer, router)

	return router.Listen(fmt.Sprintf(":%d", cfg.Port))
}

func registerSMSRoutes(appContainer app.App, router fiber.Router) {
	ctx := context.Background()
	smsUseCase := appContainer.SMSService(ctx)

	smsHandler := NewSMSHandler(smsUseCase)

	// API v1 routes
	v1 := router.Group("/api/v1")

	// SMS routes
	sms := v1.Group("/sms")
	sms.Post("/", setTraceID(), smsHandler.SendSMS)
	sms.Get("/:id", setTraceID(), smsHandler.GetSMS)

}

// customErrorHandler handles errors consistently
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
