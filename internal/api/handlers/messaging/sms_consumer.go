package messaging

import (
	"context"
	"sms/internal/usecase/sms"
)

type ConsumerHandler struct {
	smsService sms.Service

	// TODO:add logger
}

func NewSMSConsumer(smsService sms.Service) *ConsumerHandler {
	return &ConsumerHandler{
		smsService: smsService,
	}
}

func (h *ConsumerHandler) HandleDebitedSMS(ctx context.Context, message []byte) error {
	//TODO: add process sms debited event here
	return nil
}
