package models

import (
	"time"

	"github.com/MBvisti/mortenvistisen/pkg/validation"
	"github.com/google/uuid"
)

type Newsletter struct {
	ID         uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Title      string
	Content    string
	ReleasedAt time.Time
	Released   bool
}

var CreateNewsletterValidations = func() map[string][]validation.Rule {
	return map[string][]validation.Rule{
		"Title": {validation.RequiredRule},
	}
}

var ReleaseNewsletterValidations = func() map[string][]validation.Rule {
	return map[string][]validation.Rule{
		"ID":         {validation.RequiredRule},
		"Title":      {validation.RequiredRule, validation.MinLenRule(3)},
		"Content":    {validation.RequiredRule, validation.MinLenRule(3)},
		"Released":   {validation.RequiredRule, validation.MustBeTrueRule},
		"ReleasedAt": {validation.RequiredRule},
	}
}

func (n Newsletter) CanBeReleased() error {
	return validation.Validate(n, ReleaseNewsletterValidations())
}
