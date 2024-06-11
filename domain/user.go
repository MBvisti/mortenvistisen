package domain

import (
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
	if err := usr.Validate(confirmPassword); err != nil {
		return User{}, err
	}

	return usr, nil
}

var UserValidations = map[string][]Rule{
	"ID":       {RequiredRule},
	"Name":     {RequiredRule, MinLenRule(2), MaxLenRule(50)},
	"Password": {RequiredRule, MinLenRule(6), PasswordMatchConfirmRule},
	"Mail":     {RequiredRule, EmailRule},
}

func (u User) Validate(confirmPassword string) error {
	var errors []ValidationErr
	for field, rules := range UserValidations {
		switch field {
		case "ID":
			idValidationErr := ErrValidation{
				FieldName:  "ID",
				FieldValue: u.ID,
			}
			for _, rule := range rules {
				if err := checkRule(u.ID, rule); err != nil {
					idValidationErr.Violations = append(
						idValidationErr.Violations,
						err,
					)
				}
			}
			errors = append(errors, idValidationErr)
		case "Name":
			nameValidationErr := ErrValidation{
				FieldName:  "Name",
				FieldValue: u.Name,
			}
			for _, rule := range rules {
				if err := checkRule(u.Name, rule); err != nil {
					nameValidationErr.Violations = append(
						nameValidationErr.Violations,
						err,
					)
				}
			}
			errors = append(errors, nameValidationErr)
		case "Password":
			passwordValidationErrs := ErrValidation{
				FieldName:  "Password",
				FieldValue: u.Password,
			}
			for _, rule := range rules {
				compareable, ok := rule.(Compareable)
				if !ok {
					if err := checkRule(u.Password, rule); err != nil {
						passwordValidationErrs.Violations = append(
							passwordValidationErrs.Violations,
							err,
						)
					}
				}
				if ok {
					err := checkComparableRule(u.Password, confirmPassword, rule, compareable)
					passwordValidationErrs.Violations = append(
						passwordValidationErrs.Violations,
						err,
					)
				}
			}
			errors = append(errors, passwordValidationErrs)
		case "Mail":
			emailValidationErr := ErrValidation{
				FieldName:  "Mail",
				FieldValue: u.Mail,
			}
			for _, rule := range rules {
				if err := checkRule(u.Mail, rule); err != nil {
					emailValidationErr.Violations = append(
						emailValidationErr.Violations,
						err,
					)
				}
			}
			errors = append(errors, emailValidationErr)
		}
	}

	e := constructValidationErrors(errors...)
	if len(e) > 0 {
		return e
	}

	return nil
}
