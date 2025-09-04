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
	smsUseCase *sms.UseCase
}

func NewSMSHandler(smsUseCase *sms.UseCase) *SMSHandler {
	return &SMSHandler{
		smsUseCase: smsUseCase,
	}
}

func (h *SMSHandler) SendSMS(c *fiber.Ctx) error {
	var req dto.SendSMSRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	if req.Content == "" || req.Receiver == "" || req.UserID == "" {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "validation_error",
			Message: "content, receiver, and user_id are required",
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

func (h *SMSHandler) GetSMS(c *fiber.Ctx) error {
	smsID := c.Params("id")
	if smsID == "" {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "SMS ID is required",
		})
	}

	filter := smsdomain.Filter{
		ID: &smsID,
	}

	ctx := c.Context()
	smsMessage, err := h.smsUseCase.GetSMSByID(ctx, filter)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(dto.ErrorResponse{
			Error:   "not_found",
			Message: "SMS not found",
		})
	}

	response := dto.GetSMSResponse{
		ID:        smsMessage.ID,
		UserID:    smsMessage.UserID,
		Content:   smsMessage.Content,
		Receiver:  smsMessage.Receiver,
		Provider:  smsMessage.Provider,
		Status:    string(smsMessage.Status),
		CreatedAt: smsMessage.CreatedAt,
		UpdatedAt: smsMessage.UpdatedAt,
	}

	if !smsMessage.DeliveredAt.IsZero() {
		response.DeliveredAt = &smsMessage.DeliveredAt
	}

	if smsMessage.FailureCode != "" {
		response.FailureCode = smsMessage.FailureCode
	}

	return c.JSON(response)
}

func (h *SMSHandler) GetUserSMS(c *fiber.Ctx) error {
	userID := c.Params("userID")
	if userID == "" {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "User ID is required",
		})
	}

	filter := smsdomain.Filter{
		UserID: &userID,
	}

	ctx := c.Context()
	smsMessage, err := h.smsUseCase.GetSMSByID(ctx, filter)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(dto.ErrorResponse{
			Error:   "not_found",
			Message: "No SMS found for this user",
		})
	}

	response := dto.GetSMSResponse{
		ID:        smsMessage.ID,
		UserID:    smsMessage.UserID,
		Content:   smsMessage.Content,
		Receiver:  smsMessage.Receiver,
		Provider:  smsMessage.Provider,
		Status:    string(smsMessage.Status),
		CreatedAt: smsMessage.CreatedAt,
		UpdatedAt: smsMessage.UpdatedAt,
	}

	if !smsMessage.DeliveredAt.IsZero() {
		response.DeliveredAt = &smsMessage.DeliveredAt
	}

	if smsMessage.FailureCode != "" {
		response.FailureCode = smsMessage.FailureCode
	}

	return c.JSON(response)
}
