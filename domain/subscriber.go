package domain

import (
	"reflect"
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

var BuildSubscriberValidations = func() map[string][]Rule {
	return map[string][]Rule{
		"ID":    {RequiredRule},
		"Email": {RequiredRule, EmailRule},
		// "SubscribedAt": {RequiredRule},
		"Referer": {RequiredRule},
	}
}

func (s Subscriber) Validate(validations map[string][]Rule) error {
	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)
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

	if err := sub.Validate(BuildSubscriberValidations()); err != nil {
		return Subscriber{}, err
	}

	return sub, nil
}
