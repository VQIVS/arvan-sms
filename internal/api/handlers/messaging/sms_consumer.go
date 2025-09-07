package messaging

import (
	"context"
	"encoding/json"
	"sms/config"
	smsDomain "sms/internal/domain/sms"
	"sms/internal/usecase/sms"
	"sms/pkg/logger"
	"sms/pkg/rabbit"
)

type ConsumerHandler struct {
	smsService sms.Service
	log        *logger.Logger
	consumer   *rabbit.Consumer
	config     config.Config
}

func NewSMSConsumer(smsService sms.Service, log *logger.Logger, rabbitConn *rabbit.RabbitConn, cfg config.Config) *ConsumerHandler {
	return &ConsumerHandler{
		smsService: smsService,
		log:        log,
		consumer:   rabbit.NewConsumer(rabbitConn),
		config:     cfg,
	}
}

func (h *ConsumerHandler) HandleDebitedSMS(ctx context.Context, message []byte) error {
	h.log.Info(ctx, "received billing completed message", "message_size", len(message))

	var msg smsDomain.SMSBillingCompleted
	err := json.Unmarshal(message, &msg)
	if err != nil {
		h.log.Error(ctx, "failed to unmarshal billing completed message", "error", err, "raw_message", string(message))
		return err
	}

	h.log.Info(ctx, "processing billing completed event", "sms_id", msg.SMSID, "transaction_id", msg.TransactionID, "user_id", msg.UserID)

	err = h.smsService.ProcessDebitedSMS(ctx, msg)
	if err != nil {
		h.log.Error(ctx, "failed to process billing completed event", "error", err, "sms_id", msg.SMSID, "transaction_id", msg.TransactionID)
		return err
	}

	h.log.Info(ctx, "billing completed event processed successfully", "sms_id", msg.SMSID, "transaction_id", msg.TransactionID)
	return nil
}

func (h *ConsumerHandler) Run(ctx context.Context) error {
	h.log.Info(ctx, "initializing SMS consumer")

	if err := h.consumer.SetQos(1); err != nil {
		h.log.Error(ctx, "failed to set consumer QoS", "error", err)
		return err
	}
	h.log.Info(ctx, "consumer QoS set successfully", "prefetch_count", 1)

	for _, queue := range h.config.RabbitMQ.Queues {
		switch queue.Name {
		//TODO: change to correct queue name and do not hardcode here
		case rabbit.SMSBillingCompletedQueue:
			h.consumer.Subscribe(queue.Name, func(message []byte) error {
				return h.HandleDebitedSMS(ctx, message)
			})
			h.log.Info(ctx, "subscribed to queue successfully", "queue", queue.Name, "routing_key", queue.RoutingKey)
		default:
			h.log.Info(ctx, "skipping unknown queue in configuration", "queue", queue.Name)
		}
	}

	h.log.Info(ctx, "starting SMS consumer workers")
	if err := h.consumer.StartConsume(); err != nil {
		h.log.Error(ctx, "failed to start consumer workers", "error", err)
		return err
	}
	h.log.Info(ctx, "SMS consumer workers started successfully")

	<-ctx.Done()
	h.log.Info(ctx, "SMS consumer shutdown signal received")
	return ctx.Err()
}
