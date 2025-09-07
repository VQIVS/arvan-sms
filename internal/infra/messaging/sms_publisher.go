package messaging

import (
	"context"
	"fmt"
	"sms/internal/domain/sms"
	"sms/pkg/logger"
	"sms/pkg/rabbit"
)

const (
	billingRequestedRoutingKey = "sms.billing.requested"
	billingRefundedRoutingKey  = "sms.billing.refunded"
	exchange                   = "amq.topic"
)

type SMSPublisher struct {
	publisher *rabbit.Publisher
	log       *logger.Logger
}

func NewSMSPublisher(conn *rabbit.RabbitConn, log *logger.Logger) sms.EventPublisher {
	return &SMSPublisher{
		publisher: rabbit.NewPublisher(conn),
		log:       log,
	}
}

func (p *SMSPublisher) PublishEvent(ctx context.Context, event sms.DomainEvent) error {
	switch event.EventType() {
	case sms.EventTypeBillingRequested:
		p.log.Info(ctx, "publishing billing requested event", "sms_id", event.AggregateID(), "routing_key", billingRequestedRoutingKey)
		return p.publisher.Publish(billingRequestedRoutingKey, exchange, event)
	case sms.EventTypeBillingRefunded:
		p.log.Info(ctx, "publishing billing refunded event", "transaction_id", event.AggregateID(), "routing_key", billingRefundedRoutingKey)
		return p.publisher.Publish(billingRefundedRoutingKey, exchange, event)
	default:
		return fmt.Errorf("unknown event type: %s", event.EventType())
	}
}
