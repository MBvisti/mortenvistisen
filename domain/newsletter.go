package domain

import (
	"time"

	"github.com/go-playground/validator/v10"
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

var NewsletterValidations = map[string][]Rule{
	"ID":          {RequiredRule},
	"Title":       {RequiredRule, MinLenRule(3)},
	"Edition":     {RequiredRule},
	"Paragraphs":  {RequiredRule, MinLenRule(1)},
	"ArticleSlug": {RequiredRule},
	"ReleasedAt":  {RequiredRule},
	"Released  ":  {RequiredRule},
}

func (n Newsletter) Validate() error {
	var errors []ValidationErr
	for field, rules := range NewsletterValidations {
		switch field {
		case "ID":
			idValidationErr := ErrValidation{
				FieldName:  "ID",
				FieldValue: n.ID,
			}
			for _, rule := range rules {
				if err := checkRule(n.ID, rule); err != nil {
					idValidationErr.Violations = append(
						idValidationErr.Violations,
						err,
					)
				}
			}
			errors = append(errors, idValidationErr)
		case "Title":
			titleValidationErr := ErrValidation{
				FieldName:  "Title",
				FieldValue: n.Title,
			}
			for _, rule := range rules {
				if err := checkRule(n.Title, rule); err != nil {
					titleValidationErr.Violations = append(
						titleValidationErr.Violations,
						err,
					)
				}
			}
			errors = append(errors, titleValidationErr)
		case "Edition":
			editionValidationErr := ErrValidation{
				FieldName:  "Edition",
				FieldValue: n.Edition,
			}
			for _, rule := range rules {
				if err := checkRule(n.Edition, rule); err != nil {
					editionValidationErr.Violations = append(
						editionValidationErr.Violations,
						err,
					)
				}
			}
			errors = append(errors, editionValidationErr)
		case "Paragraphs":
			paragraphsValidationErr := ErrValidation{
				FieldName:  "Paragraphs",
				FieldValue: n.Paragraphs,
			}
			for _, rule := range rules {
				if err := checkRule(n.Paragraphs, rule); err != nil {
					paragraphsValidationErr.Violations = append(
						paragraphsValidationErr.Violations,
						err,
					)
				}
			}
			errors = append(errors, paragraphsValidationErr)
		case "ArticleSlug":
			articleValidationErr := ErrValidation{
				FieldName:  "ArticleSlug",
				FieldValue: n.ArticleSlug,
			}
			for _, rule := range rules {
				if err := checkRule(n.ArticleSlug, rule); err != nil {
					articleValidationErr.Violations = append(
						articleValidationErr.Violations,
						err,
					)
				}
			}
			errors = append(errors, articleValidationErr)
		case "ReleasedAt":
			releasedAtValidationErr := ErrValidation{
				FieldName:  "ReleasedAt",
				FieldValue: n.ReleasedAt,
			}
			for _, rule := range rules {
				if err := checkRule(n.ReleasedAt, rule); err != nil {
					releasedAtValidationErr.Violations = append(
						releasedAtValidationErr.Violations,
						err,
					)
				}
			}
			errors = append(errors, releasedAtValidationErr)
		case "Released":
			releasedValidationErr := ErrValidation{
				FieldName:  "Released",
				FieldValue: n.Released,
			}
			for _, rule := range rules {
				if err := checkRule(n.Released, rule); err != nil {
					releasedValidationErr.Violations = append(
						releasedValidationErr.Violations,
						err,
					)
				}
			}
			errors = append(errors, releasedValidationErr)
		}
	}

	e := constructValidationErrors(errors...)
	if len(e) > 0 {
		return e
	}

	return nil
}

func NewNewsletter(
	title string,
	edition int32,
	releasedAt time.Time,
	released bool,
	paragraphs []string,
	articleSlug string,
) (Newsletter, error) {
	now := time.Now()

	newsletter := Newsletter{
		uuid.New(),
		now,
		now,
		title,
		edition,
		releasedAt,
		released,
		paragraphs,
		articleSlug,
	}

	if err := newsletter.Validate(); err != nil {
		return Newsletter{}, err
	}

	return newsletter, nil
}

type UpdateNewsletterPayload struct {
	ID          uuid.UUID
	Title       string   `validate:"required,gte=3"`
	Edition     int32    `validate:"required,gte=1"`
	Paragraphs  []string `validate:"required,gte=1"`
	ArticleSlug string   `validate:"required"`
}

func (n Newsletter) Update(
	payload UpdateNewsletterPayload,
	v *validator.Validate,
) (Newsletter, error) {
	if err := v.Struct(payload); err != nil {
		return Newsletter{}, err
	}

	return Newsletter{
		ID:          payload.ID,
		Title:       payload.Title,
		Edition:     payload.Edition,
		Paragraphs:  payload.Paragraphs,
		ArticleSlug: payload.ArticleSlug,
	}, nil
}

func InitilizeNewsletter(
	title string,
	edition int32,
	paragraphs []string,
	articleSlug string,
) Newsletter {
	return Newsletter{
		ID:          uuid.New(),
		Title:       title,
		Edition:     edition,
		Paragraphs:  paragraphs,
		ArticleSlug: articleSlug,
	}
}
