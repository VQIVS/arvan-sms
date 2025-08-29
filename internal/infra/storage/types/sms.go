package types

import (
	"time"
)

type SMS struct {
	Base
	Content     string
	Receiver    string
	Provider    string
	Status      string
	DeliveredAt *time.Time
	FailureCode *string
}
