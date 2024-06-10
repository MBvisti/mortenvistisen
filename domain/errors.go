package domain

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	ErrInvalidEmail      = errors.New("provided email is invalid")
	ErrInvalidUsername   = errors.New("provided username is invalid")
	ErrIsRequired        = errors.New("value is required")
	ErrTooShort          = errors.New("value is too short")
	ErrTooLong           = errors.New("value is too long")
	ErrPasswordDontMatch = errors.New("the two passwords must match")
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func isEmailValid(e string) bool {
	return emailRegex.MatchString(e)
}

type ValidationErr interface {
	Error() string
	Field() string
	Value() any
	Causes() []error
}

var baseErrMsg = "Field: '%s' with Value: '%v' has Error(s): validation failed due to '%v'"

type ValidationErrs []ValidationErr

func (ve ValidationErrs) Error() string {
	var errMsg string
	for _, err := range ve {
		errMsg += err.Error() + "; "
	}
	return errMsg
}

type ErrValidation struct {
	FieldValue any
	FieldName  string
	Violations []error
}

func (e ErrValidation) Field() string {
	return e.FieldName
}

func (e ErrValidation) Value() any {
	return e.FieldValue
}

func (e ErrValidation) Causes() []error {
	return e.Violations
}

func (e ErrValidation) Error() string {
	var causes string
	for i, violation := range e.Violations {
		if i == 0 {
			causes = violation.Error()
		}

		if i != 0 {
			causes = causes + ", " + violation.Error()
		}
	}

	return fmt.Sprintf(
		baseErrMsg,
		e.FieldName,
		e.FieldValue,
		causes,
	)
}

func constructValidationErrors(valiErrs ...ValidationErr) ValidationErrs {
	var validationErrors ValidationErrs
	for _, valiErr := range valiErrs {
		if len(valiErr.Causes()) > 0 {
			validationErrors = append(validationErrors, valiErr)
		}
	}

	return validationErrors
}
