package domain

import (
	"reflect"
	"time"

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
				errVal.ViolationsForHuman = append(
					errVal.ViolationsForHuman,
					rule.ViolationForHumans(name),
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
		time.Time{},
		false,
		paragraphs,
		articleSlug,
	}

	if err := newsletter.Validate(BuildNewsletterValidations()); err != nil {
		return Newsletter{}, err
	}

	return newsletter, nil
}

func (n Newsletter) Update(
	title string,
	edition int32,
	paragraphs []string,
	articleSlug string,
) (Newsletter, error) {
	n.Title = title
	n.Edition = edition
	n.Paragraphs = paragraphs
	n.ArticleSlug = articleSlug
	n.UpdatedAt = time.Now()

	if err := n.Validate(BuildNewsletterValidations()); err != nil {
		return Newsletter{}, err
	}

	return n, nil
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
