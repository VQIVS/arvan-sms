package http

import (
	"sms-dispatcher/pkg/context"
	"sms-dispatcher/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func setUserContext(c *fiber.Ctx) error {
	c.SetUserContext(context.NewAppContext(c.UserContext(), context.WithLogger(logger.GetLogger())))
	return c.Next()
}
