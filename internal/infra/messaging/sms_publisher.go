package messaging

import (
	"context"
	"fmt"
	"sms/internal/domain/sms"
	"sms/pkg/rabbit"
)

const exchange = "amq.topic"

// TODO: don't hardcode this here
const debitRoutingKey = "sms.debit.balance"
const refundRoutingKey = "sms.refund.balance"

type SMSPublisher struct {
	publisher *rabbit.Publisher
	//TODO: add logger
}

func NewSMSPublisher(conn *rabbit.RabbitConn) sms.Publisher {
	return &SMSPublisher{
		publisher: rabbit.NewPublisher(conn),
	}
}

func (p *SMSPublisher) PublishEvent(ctx context.Context, event sms.SMSEvent) error {
	switch event.EventType() {
	case sms.EventTypeDebit:
		return p.publisher.Publish(debitRoutingKey, exchange, event)
	case sms.EventTypeRefund:
		return p.publisher.Publish(refundRoutingKey, exchange, event)
	default:
		return fmt.Errorf("unknown event type: %s", event.EventType())
	}
}
