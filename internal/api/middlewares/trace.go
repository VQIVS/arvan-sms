package middlewares

import (
	"sms/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func FiberTraceMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := logger.WithTraceID(c.Context())
		c.SetUserContext(ctx)

		traceID := logger.GetTraceID(ctx)
		c.Set("X-Trace-ID", traceID)
		return c.Next()
	}
}
