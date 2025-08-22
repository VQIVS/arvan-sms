package constants

const (
	// consumer routing keys
	KeySMSUpdate = "finance.balance.update"
	// producer routing keys
	KeyBalanceUpdate = "sms.status.update"

	TopicExchange = "amq.topic"
	ServiceName   = "sms"
)
