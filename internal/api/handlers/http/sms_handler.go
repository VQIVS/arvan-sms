package http

import (
	"net/http"
	"sms/internal/api/dto"
	smsdomain "sms/internal/domain/sms"
	"sms/internal/usecase/sms"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SMSHandler struct {
	smsUseCase *sms.Service
}

func NewSMSHandler(smsUseCase *sms.Service) *SMSHandler {
	return &SMSHandler{
		smsUseCase: smsUseCase,
	}
}

// SendSMS godoc
// @Summary Send an SMS message
// @Description Send an SMS message to a specified receiver
// @Tags SMS
// @Accept json
// @Produce json
// @Param sms body dto.SendSMSRequest true "SMS request payload"
// @Success 201 {object} dto.SendSMSResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /sms [post]
func (h *SMSHandler) SendSMS(c *fiber.Ctx) error {
	var req dto.SendSMSRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	smsMessage := &smsdomain.SMSMessage{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Content:   req.Content,
		Receiver:  req.Receiver,
		Status:    smsdomain.SMSStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	ctx := c.Context()
	if err := h.smsUseCase.ProcessSMS(ctx, smsMessage); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "processing_error",
			Message: "Failed to process SMS",
		})
	}

	return c.Status(http.StatusCreated).JSON(dto.SendSMSResponse{
		ID:        smsMessage.ID,
		Status:    string(smsMessage.Status),
		CreatedAt: smsMessage.CreatedAt,
		Message:   "SMS queued for processing",
	})
}
