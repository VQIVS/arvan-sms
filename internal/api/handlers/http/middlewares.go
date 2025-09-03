package http

import (
	"sms/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

// TODO: make this private
func SetTraceID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := logger.WithTraceID(c.Context())
		c.SetUserContext(ctx)

		traceID := logger.GetTraceID(ctx)
		c.Set("X-Trace-ID", traceID)
		return c.Next()
	}
}
