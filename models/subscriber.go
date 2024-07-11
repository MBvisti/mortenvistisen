package models

import (
	"time"

	"github.com/MBvisti/mortenvistisen/pkg/validation"
	"github.com/google/uuid"
)

type Subscriber struct {
	ID           uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Email        string
	SubscribedAt time.Time
	Referer      string
	IsVerified   bool
}

var CreateSubscriberValidations = func() map[string][]validation.Rule {
	return map[string][]validation.Rule{
		"ID":      {validation.RequiredRule},
		"Email":   {validation.RequiredRule, validation.EmailRule},
		"Referer": {validation.RequiredRule},
	}
}
