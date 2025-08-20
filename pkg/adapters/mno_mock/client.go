package mno

import "sms-dispatcher/pkg/logger"

type StatusCode uint

const (
	SuccessCode StatusCode = 200
	FailedCode  StatusCode = 500
)

var callCount int

// Fail on every 4th call, succeed otherwise
func SendSMSViaMNO() (StatusCode, error) {
	callCount++

	if callCount%4 == 1 {
		return FailedCode, nil
	}
	return SuccessCode, nil
}

func FailAll() (StatusCode, error) {
	logger.GetLogger().Info("MNO: failing all requests =================================================")
	return FailedCode, nil
}
