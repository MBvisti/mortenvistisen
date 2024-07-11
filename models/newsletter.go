package models

import (
	"time"

	"github.com/MBvisti/mortenvistisen/pkg/validation"
	"github.com/google/uuid"
)

type Newsletter struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Edition     int32
	ReleasedAt  time.Time
	Released    bool
	Paragraphs  []string
	ArticleSlug string
}

var CreateNewsletterValidations = func() map[string][]validation.Rule {
	return map[string][]validation.Rule{
		"ArticleSlug": {validation.RequiredRule},
	}
}

var ReleaseNewsletterValidations = func() map[string][]validation.Rule {
	return map[string][]validation.Rule{
		"ID":          {validation.RequiredRule},
		"Title":       {validation.RequiredRule, validation.MinLenRule(3)},
		"Edition":     {validation.RequiredRule},
		"Paragraphs":  {validation.RequiredRule, validation.MinLenRule(1)},
		"ArticleSlug": {validation.RequiredRule},
		"Released":    {validation.RequiredRule, validation.MustBeTrueRule},
		"ReleasedAt":  {validation.RequiredRule},
	}
}

func (n Newsletter) CanBeReleased() error {
	return validation.Validate(n, ReleaseNewsletterValidations())
}
