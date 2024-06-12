package domain

import (
	"reflect"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Name           string
	Mail           string
	MailVerifiedAt time.Time
	Password       string
}

func NewUser(name, mail, password, confirmPassword string) (User, error) {
	now := time.Now()
	usr := User{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
		Mail:      mail,
		Password:  password,
	}
	if err := usr.Validate(BuildUserValidations(confirmPassword)); err != nil {
		return User{}, err
	}

	return usr, nil
}

var BuildUserValidations = func(confirm string) map[string][]Rule {
	return map[string][]Rule{
		"ID":        {RequiredRule},
		"Name":      {RequiredRule, MinLenRule(2), MaxLenRule(25)},
		"Password":  {RequiredRule, MinLenRule(6), PasswordMatchConfirmRule(confirm)},
		"Mail":      {RequiredRule, EmailRule},
		"CreatedAt": {RequiredRule},
	}
}

func (u User) Validate(validations map[string][]Rule) error {
	val := reflect.ValueOf(u)
	typ := reflect.TypeOf(u)
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

func (u User) Update(
	name string,
	mail string,
	mailVerifiedAt time.Time,
	password string,
	confirmPassword string,
) (User, error) {
	updatedUser := User{
		ID:             u.ID,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      time.Now(),
		Name:           name,
		Mail:           mail,
		MailVerifiedAt: mailVerifiedAt,
		Password:       password,
	}

	if err := updatedUser.Validate(BuildUserValidations(confirmPassword)); err != nil {
		return User{}, err
	}

	return updatedUser, nil
}
