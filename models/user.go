package models

import (
	"time"

	"github.com/MBvisti/mortenvistisen/pkg/validation"
	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Name           string
	Mail           string
	MailVerifiedAt time.Time
	Password       string
}

var CreateUserValidations = func(confirm string) map[string][]validation.Rule {
	return map[string][]validation.Rule{
		"ID":   {validation.RequiredRule},
		"Name": {validation.RequiredRule, validation.MinLenRule(2), validation.MaxLenRule(25)},
		"Password": {
			validation.RequiredRule,
			validation.MinLenRule(6),
			validation.PasswordMatchConfirmRule(confirm),
		},
		"Mail":      {validation.RequiredRule, validation.EmailRule},
		"CreatedAt": {validation.RequiredRule},
	}
}
