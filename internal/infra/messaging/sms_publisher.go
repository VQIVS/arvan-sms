package messaging

import (
	"context"
	"fmt"
	"sms/internal/domain/sms"
	"sms/pkg/rabbit"
)

const (
	billingRequestedRoutingKey = "sms.billing.requested"
	billingRefundedRoutingKey  = "sms.billing.refunded"
	exchange                   = "amq.topic"
)

type SMSPublisher struct {
	publisher *rabbit.Publisher
	//TODO: add logger
}

func NewSMSPublisher(conn *rabbit.RabbitConn) sms.EventPublisher {
	return &SMSPublisher{
		publisher: rabbit.NewPublisher(conn),
	}
}

func (p *SMSPublisher) PublishEvent(ctx context.Context, event sms.DomainEvent) error {
	switch event.EventType() {
	case sms.EventTypeBillingRequested:
		return p.publisher.Publish(billingRequestedRoutingKey, exchange, event)
	case sms.EventTypeBillingRefunded:
		return p.publisher.Publish(billingRefundedRoutingKey, exchange, event)
	default:
		return fmt.Errorf("unknown event type: %s", event.EventType())
	}
}
