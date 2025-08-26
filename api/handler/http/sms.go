package http

import (
	"sms-dispatcher/api/presenter"
	"sms-dispatcher/api/service"
	"strconv"

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
func SendSMSMessage(svc *service.SMSService) fiber.Handler {
	return func(c *fiber.Ctx) error {
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

// GetSMSMessage godoc
// @Summary Get SMS message
// @Description Get a SMS message by ID
// @Tags SMS
// @Accept json
// @Produce json
// @Param id path int true "SMS ID"
// @Success 200 {object} presenter.SMSResp
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /sms/{id} [get]
func GetSMSMessage(svc *service.SMSService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		smsIDStr := c.Params("id")
		smsIDUint, err := strconv.ParseUint(smsIDStr, 10, 32)
		if err != nil {
			return fiber.ErrBadRequest
		}
		resp, err := svc.GetSMSMessage(c.UserContext(), uint(smsIDUint))
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}
