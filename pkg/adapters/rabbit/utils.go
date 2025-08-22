package rabbit

import (
	"sms-dispatcher/pkg/constants"
)

func GetQueueName(key string) string {
	return constants.ServiceName + "_" + key
}
