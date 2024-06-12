package domain

import (
	"reflect"
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

var BuildNewsletterValidations = func() map[string][]Rule {
	return map[string][]Rule{
		"ID":          {RequiredRule},
		"Title":       {RequiredRule, MinLenRule(3)},
		"Edition":     {RequiredRule},
		"Paragraphs":  {RequiredRule, MinLenRule(1)},
		"ArticleSlug": {RequiredRule},
		"ReleasedAt":  {RequiredRule},
		"Released  ":  {RequiredRule},
	}
}

func (n Newsletter) Validate(validations map[string][]Rule) error {
	val := reflect.ValueOf(n)
	typ := reflect.TypeOf(n)
	var errors []ValidationErr
	for i := 0; i < val.NumField(); i++ {
		value := val.Field(i)
		name := typ.Field(i).Name

		errVal := ErrValidation{
			FieldValue: value,
			FieldName:  name,
		}

		for _, rule := range validations[name] {
			if rule.IsViolated(GetFieldValue(value)) {
				errVal.Violations = append(
					errVal.Violations,
					rule.Violation(),
				)
			}
		}

		if len(errVal.Violations) > 0 {
			errors = append(errors, errVal)
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

	if err := newsletter.Validate(BuildNewsletterValidations()); err != nil {
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
