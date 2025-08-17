package types

import "gorm.io/gorm"

type Status string

const (
	Pending   Status = "pending"
	Delivered Status = "delivered"
	Failed    Status = "failed"
)

type SMS struct {
	gorm.Model
	Recipient string `gorm:"not null;index"`
	Message   string `gorm:"not null"`
	Status    string `gorm:"not null;index"`
}
