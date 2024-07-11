package models

import (
	"github.com/MBvisti/mortenvistisen/pkg/validation"
	"github.com/google/uuid"
)

type Tag struct {
	ID   uuid.UUID
	Name string
}

var CreateTagValidations = func() map[string][]validation.Rule {
	return map[string][]validation.Rule{
		"ID":   {validation.RequiredRule},
		"Name": {validation.RequiredRule, validation.MinLenRule(2)},
	}
}
