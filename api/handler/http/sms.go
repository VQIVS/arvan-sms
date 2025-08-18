package http

import (
	"sms-dispatcher/api/presenter"
	"sms-dispatcher/api/service"

	"github.com/gofiber/fiber/v2"
)

// SendSMSMessage godoc
// @Summary Send SMS message
// @Description Send a SMS message
// @Tags SMS
// @Accept json
// @Produce json
// @Param request body presenter.SendSMSReq true "Send SMS request"
// @Success 200 {object} presenter.SendSMSResp
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /sms/send [post]
func SendSMSMessage(svcGetter ServiceGetter[*service.SMSService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		svc := svcGetter(c.UserContext())
		var req presenter.SendSMSReq
		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}
		resp, err := svc.SendSMS(c.UserContext(), &req)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}
