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

	ctx := c.UserContext()
	if err := h.smsUseCase.CreateAndBillSMS(ctx, smsMessage); err != nil {
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

// GetSMSByID godoc
// @Summary Get an SMS message by ID
// @Description Retrieve SMS message details by its ID
// @Tags SMS
// @Accept json
// @Produce json
// @Param id path string true "SMS ID"
// @Success 200 {object} dto.GetSMSResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /sms/{id} [get]
func (h *SMSHandler) GetSMSByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "SMS ID is required",
		})
	}

	ctx := c.UserContext()
	smsMessage, err := h.smsUseCase.GetSMSByID(ctx, smsdomain.Filter{ID: &id})
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(dto.ErrorResponse{
			Error:   "not_found",
			Message: "SMS not found",
		})
	}

	var deliveredAt *time.Time
	if !smsMessage.DeliveredAt.IsZero() {
		deliveredAt = &smsMessage.DeliveredAt
	}

	return c.Status(http.StatusOK).JSON(dto.GetSMSResponse{
		ID:          smsMessage.ID,
		UserID:      smsMessage.UserID,
		Content:     smsMessage.Content,
		Receiver:    smsMessage.Receiver,
		Provider:    smsMessage.Provider,
		Status:      string(smsMessage.Status),
		DeliveredAt: deliveredAt,
		FailureCode: smsMessage.FailureCode,
		CreatedAt:   smsMessage.CreatedAt,
		UpdatedAt:   smsMessage.UpdatedAt,
	})
}
