package mapper

import (
	"sms/internal/domain/sms"
	"sms/internal/infra/storage/types"
)

func TODomain(model types.SMS) *sms.SMSMessage {
	return &sms.SMSMessage{
		ID:          model.ID,
		UserID:      model.UserID,
		Content:     model.Content,
		Receiver:    model.Receiver,
		Provider:    *model.Provider,
		Status:      sms.SMSStatus(model.Status),
		DeliveredAt: *model.DeliveredAt,
		FailureCode: *model.FailureCode,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
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
