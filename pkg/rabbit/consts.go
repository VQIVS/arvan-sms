package rabbit

const (
	// consumer will use this queue
	SMSBillingCompletedQueue = "sms_billing.debit.completed"
	// producer will use this routing key to publish billing requested event
	BillingRequestedRoutingKey = "billing.debit.request"
	BillingRefundedRoutingKey  = "billing.refund.request"
	Exchange                   = "amq.topic"
)
