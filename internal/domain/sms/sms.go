package sms

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Repo interface {
	GetByFilter(ctx context.Context, filter Filter) (*SMSMessage, error)
	Create(ctx context.Context, message *SMSMessage) error
	Update(ctx context.Context, ID string, message *SMSMessage) error
	WithTx(tx *gorm.DB) Repo
}

type SMSStatus string

const (
	SMSStatusPending   SMSStatus = "pending"
	SMSStatusDelivered SMSStatus = "delivered"
	SMSStatusFailed    SMSStatus = "failed"
)

type SMSMessage struct {
	ID          string
	UserID      string
	Content     string
	Receiver    string
	Provider    string
	Status      SMSStatus
	DeliveredAt time.Time
	FailureCode string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

type Filter struct {
	ID     *string
	Status *SMSStatus
	UserID *string
}

func (s *SMSMessage) MarkAsSent(provider string) {
	s.Status = SMSStatusDelivered
	s.Provider = provider
	now := time.Now()
	s.DeliveredAt = now
	s.UpdatedAt = now
}

func (s *SMSMessage) MarkAsFailed(provider string, code string) {
	s.Status = SMSStatusFailed
	s.Provider = provider
	s.FailureCode = code
	s.UpdatedAt = time.Now()
}
