package types

import (
	"time"
)

type SMS struct {
	Base
	UserID      string
	Content     string
	Receiver    string
	Provider    *string
	Status      string
	DeliveredAt *time.Time
	FailureCode *string
}
