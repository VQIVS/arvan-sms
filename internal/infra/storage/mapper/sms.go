package mapper

import (
	"sms/internal/domain/sms"
	"sms/internal/infra/storage/types"
)

func TODomain(model types.SMS) *sms.SMSMessage {
	result := &sms.SMSMessage{
		ID:        model.ID,
		UserID:    model.UserID,
		Content:   model.Content,
		Receiver:  model.Receiver,
		Status:    sms.SMSStatus(model.Status),
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}

	// Handle nullable fields safely
	if model.Provider != nil {
		result.Provider = *model.Provider
	}

	if model.DeliveredAt != nil {
		result.DeliveredAt = *model.DeliveredAt
	}

	if model.FailureCode != nil {
		result.FailureCode = *model.FailureCode
	}

	if model.DeletedAt != nil {
		result.DeletedAt = *model.DeletedAt
	}

	return result
}

func TOStorage(sms sms.SMSMessage) *types.SMS {
	return &types.SMS{
		Base: types.Base{
			ID:        sms.ID,
			CreatedAt: sms.CreatedAt,
			UpdatedAt: sms.UpdatedAt,
			DeletedAt: &sms.DeletedAt,
		},
		UserID:      sms.UserID,
		Content:     sms.Content,
		Receiver:    sms.Receiver,
		Provider:    &sms.Provider,
		Status:      string(sms.Status),
		DeliveredAt: &sms.DeliveredAt,
		FailureCode: &sms.FailureCode,
	}
}
