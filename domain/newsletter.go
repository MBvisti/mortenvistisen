package domain

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ValidationErrsMap map[string]string

type Newsletter struct {
	ID          uuid.UUID `validate:"required"`
	Title       string    `validate:"required,gte=3"`
	Edition     int32     `validate:"required,gte=1"`
	ReleasedAt  time.Time
	Released    bool
	Paragraphs  []string `validate:"required,gte=1"`
	ArticleSlug string   `validate:"required"`
}

// CanBeReleased checks if a newsletter is ready to be released and updates the 'Released' and 'ReleasedAt' properties
func (n Newsletter) CanBeReleased(v *validator.Validate) (Newsletter, ValidationErrsMap, error) {
	if err := v.Struct(n); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return Newsletter{}, nil, errors.New("could not convert errors to ValidationErrors")
		}

		errors := make(ValidationErrsMap, len(validationErrors))
		for _, valiErr := range validationErrors {
			switch valiErr.Field() {
			case "ID":
				errors[valiErr.Field()] = "a valid uuid v4 must be provided"
			case "Title":
				errors[valiErr.Field()] = "title cannot be empty"
			case "Paragraphs":
				errors[valiErr.Field()] = "atleast one paragraph is needed"
			case "ArticleSlug":
				errors[valiErr.Field()] = "an article slug must be provided"
			case "Edition":
				errors[valiErr.Field()] = "edition is required and must be > 0"
			}
		}

		return Newsletter{}, errors, nil
	}

	n.Released = true
	n.ReleasedAt = time.Now()

	return n, nil, nil
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
