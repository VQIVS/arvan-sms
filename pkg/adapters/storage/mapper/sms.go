package mapper

import (
	"sms-dispatcher/internal/sms/domain"
	"sms-dispatcher/pkg/adapters/storage/types"

	"gorm.io/gorm"
)

func SMSDomain2Storage(sms domain.SMS) *types.SMS {
	return &types.SMS{
		Model: gorm.Model{
			ID:        uint(sms.ID),
			CreatedAt: sms.CreatedAt,
		},
		Recipient: sms.Recipient,
		Message:   sms.Message,
		Status:    string(sms.Status),
	}

}

func SMSStorage2Domain(sms *types.SMS) *domain.SMS {
	return &domain.SMS{
		ID:        domain.SMSID(sms.ID),
		Recipient: sms.Recipient,
		Message:   sms.Message,
		Status:    sms.Status,
		CreatedAt: sms.CreatedAt,
	}
}
