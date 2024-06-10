package domain

import (
	"time"

	"github.com/go-playground/validator/v10"
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

func (s Subscriber) Validate(validator *validator.Validate) error {
	emailValidationErrs := ErrValidation{
		FieldName:  "Email",
		FieldValue: s.Email,
	}
	if s.Email == "" {
		emailValidationErrs.Violations = append(emailValidationErrs.Violations, ErrIsRequired)
	}
	if !isEmailValid(s.Email) {
		emailValidationErrs.Violations = append(emailValidationErrs.Violations, ErrInvalidEmail)
	}

	refererValidationErrs := ErrValidation{
		FieldName:  "Referer",
		FieldValue: s.Referer,
	}
	if s.Referer == "" {
		refererValidationErrs.Violations = append(refererValidationErrs.Violations, ErrIsRequired)
	}

	subbedAtValidationErrs := ErrValidation{
		FieldName:  "SubscribedAt",
		FieldValue: s.SubscribedAt,
	}
	if s.SubscribedAt.IsZero() {
		subbedAtValidationErrs.Violations = append(emailValidationErrs.Violations, ErrIsRequired)
	}

	e := constructValidationErrors(
		emailValidationErrs,
		refererValidationErrs,
		subbedAtValidationErrs,
	)
	if len(e) > 0 {
		return e
	}

	return nil
}

func NewSubscriber(
	email, referer string,
	subscribedAt time.Time,
	isVerified bool,
	validator *validator.Validate,
) (Subscriber, error) {
	now := time.Now()

	sub := Subscriber{
		uuid.New(),
		now,
		now,
		email,
		subscribedAt,
		referer,
		isVerified,
	}

	if err := sub.Validate(validator); err != nil {
		return Subscriber{}, err
	}

	return sub, nil
}
