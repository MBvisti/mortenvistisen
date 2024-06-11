package domain

import (
	"time"

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

var SubscriberValidations = map[string][]Rule{
	"ID":           {RequiredRule},
	"Email":        {RequiredRule, EmailRule},
	"SubscribedAt": {RequiredRule},
	"Referer":      {RequiredRule},
}

func (s Subscriber) Validate() error {
	var errors []ValidationErr
	for field, rules := range SubscriberValidations {
		switch field {
		case "ID":
			idValidationErr := ErrValidation{
				FieldName:  "ID",
				FieldValue: s.ID,
			}
			for _, rule := range rules {
				if err := checkRule(s.ID, rule); err != nil {
					idValidationErr.Violations = append(
						idValidationErr.Violations,
						err,
					)
				}
			}
			errors = append(errors, idValidationErr)
		case "Email":
			emailValidationErr := ErrValidation{
				FieldName:  "Email",
				FieldValue: s.Email,
			}
			for _, rule := range rules {
				if err := checkRule(s.Email, rule); err != nil {
					emailValidationErr.Violations = append(
						emailValidationErr.Violations,
						err,
					)
				}
			}
			errors = append(errors, emailValidationErr)
		case "SubscribedAt":
			subscribedAtValidationErr := ErrValidation{
				FieldName:  "SubscribedAt",
				FieldValue: s.SubscribedAt,
			}
			for _, rule := range rules {
				if err := checkRule(s.SubscribedAt, rule); err != nil {
					subscribedAtValidationErr.Violations = append(
						subscribedAtValidationErr.Violations,
						err,
					)
				}
			}
			errors = append(errors, subscribedAtValidationErr)
		case "Referer":
			refererValidationErr := ErrValidation{
				FieldName:  "Referer",
				FieldValue: s.Referer,
			}
			for _, rule := range rules {
				if err := checkRule(s.Referer, rule); err != nil {
					refererValidationErr.Violations = append(
						refererValidationErr.Violations,
						err,
					)
				}
			}
			errors = append(errors, refererValidationErr)
		}
	}

	e := constructValidationErrors(errors...)
	if len(e) > 0 {
		return e
	}

	return nil
}

func NewSubscriber(
	email string,
	referer string,
	subscribedAt time.Time,
	isVerified bool,
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

	if err := sub.Validate(); err != nil {
		return Subscriber{}, err
	}

	return sub, nil
}
