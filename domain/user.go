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

// type Validations string

// var UserValidations = map[string][]string{
// 	"Name": {"required", "gte=2"},
// }

func (u User) Validate(confirmPassword string) error {
	nameValidationErr := ErrValidation{
		FieldName:  "Name",
		FieldValue: u.Name,
	}
	if u.Name == "" {
		nameValidationErr.Violations = append(nameValidationErr.Violations, ErrIsRequired)
	}
	if len(u.Name) < 2 {
		nameValidationErr.Violations = append(nameValidationErr.Violations, ErrTooShort)
	}

	passwordValidationErrs := ErrValidation{
		FieldName: "Password",
	}
	if u.Password == "" {
		passwordValidationErrs.Violations = append(passwordValidationErrs.Violations, ErrIsRequired)
	}
	if u.Password != confirmPassword {
		passwordValidationErrs.Violations = append(
			passwordValidationErrs.Violations,
			ErrPasswordDontMatch,
		)
	}

	emailValidationErrs := ErrValidation{
		FieldName:  "Mail",
		FieldValue: u.Mail,
	}
	if u.Mail == "" {
		emailValidationErrs.Violations = append(
			emailValidationErrs.Violations,
			ErrIsRequired,
		)
	}
	if !isEmailValid(u.Mail) {
		emailValidationErrs.Violations = append(
			emailValidationErrs.Violations,
			ErrInvalidEmail,
		)
	}

	e := constructValidationErrors(
		emailValidationErrs,
		nameValidationErr,
		passwordValidationErrs,
	)
	if len(e) > 0 {
		return e
	}

	return nil
}
