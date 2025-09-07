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
	var msg smsDomain.SMSBillingCompleted
	err := json.Unmarshal(message, &msg)
	if err != nil {
		h.log.Info(ctx, "failed to unmarshal message", "error", err)
		return err
	}
	err = h.smsService.ProcessDebitedSMS(ctx, msg)
	if err != nil {
		h.log.Info(ctx, "failed to process debited sms", "error", err, "sms_id", msg.SMSID)
		return err
	}
	h.log.Logger.Info("successfully processed debited sms", "sms_id", msg.SMSID)
	return nil
}

func (h *ConsumerHandler) Run(ctx context.Context) error {
	if err := h.consumer.SetQos(1); err != nil {
		h.log.Info(ctx, "failed to set consumer QoS", "error", err)
		return err
	}

	for _, queue := range h.config.RabbitMQ.Queues {
		switch queue.Name {
		//TODO: change to correct queue name and do not hardcode here
		case rabbit.SMSBillingCompletedQueue:
			h.consumer.Subscribe(queue.Name, func(message []byte) error {
				return h.HandleDebitedSMS(ctx, message)
			})
			h.log.Logger.Info("subscribed to queue", "queue", queue.Name)
		default:
			h.log.Logger.Warn("unknown queue in configuration", "queue", queue.Name)
		}
	}

	h.log.Logger.Info("starting SMS consumer")
	if err := h.consumer.StartConsume(); err != nil {
		h.log.Info(ctx, "failed to start consumer", "error", err)
		return err
	}

	<-ctx.Done()
	h.log.Logger.Info("SMS consumer stopped")
	return ctx.Err()
}
