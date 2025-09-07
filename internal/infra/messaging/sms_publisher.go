package messaging

import (
	"context"
	"fmt"
	"sms/internal/domain/sms"
	"sms/pkg/logger"
	"sms/pkg/rabbit"
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
		p.log.Info(ctx, "publishing billing requested event", "sms_id", event.AggregateID(), "routing_key", rabbit.BillingRequestedRoutingKey)
		return p.publisher.Publish(rabbit.BillingRequestedRoutingKey, rabbit.Exchange, event)
	case sms.EventTypeBillingRefunded:
		p.log.Info(ctx, "publishing billing refunded event", "transaction_id", event.AggregateID(), "routing_key", rabbit.BillingRefundedRoutingKey)
		return p.publisher.Publish(rabbit.BillingRefundedRoutingKey, rabbit.Exchange, event)
	default:
		return fmt.Errorf("unknown event type: %s", event.EventType())
	}
}
