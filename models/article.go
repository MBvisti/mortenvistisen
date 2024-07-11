package models

import (
	"time"

	"github.com/MBvisti/mortenvistisen/pkg/validation"
	"github.com/google/uuid"
)

type Article struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	HeaderTitle string
	Filename    string
	Slug        string
	Excerpt     string
	Draft       bool
	ReleaseDate time.Time
	ReadTime    int32
	Tags        []Tag
}

var CreateArticleValidations = func() map[string][]validation.Rule {
	return map[string][]validation.Rule{
		"ID":          {validation.RequiredRule},
		"Title":       {validation.RequiredRule, validation.MinLenRule(2)},
		"HeaderTitle": {validation.RequiredRule, validation.MinLenRule(2)},
		"Excerpt": {
			validation.RequiredRule,
			validation.MinLenRule(130),
			validation.MaxLenRule(160),
		},
		"ReadTime": {validation.RequiredRule},
		"Filename": {validation.RequiredRule},
	}
}
