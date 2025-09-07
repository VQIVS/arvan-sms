package external

import (
	"sms/internal/domain/sms"
)

func DefaultSMSProvider() sms.SMSProvider {
	return RandomFailSMSProvider(0.1)
}
