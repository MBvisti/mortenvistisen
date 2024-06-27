package models

import (
	"errors"
	"log/slog"

	"github.com/go-playground/validator/v10"
)

var (
	ErrNoRowWithIdentifier = errors.New("could not find requested row in database")
	ErrNewsletterNotFound  = errors.New("could not find a newsletter for the provided identifier")
	ErrUnrecoverableEvent  = errors.New("an error occurred that could not be recovered from")
	ErrSubscriberExists    = errors.New(
		"there is already a subscriber registered with the provided data",
	)
)

type ErrValidation struct {
	ValiErr error
}

func (e ErrValidation) Error() string {
	return "domain object could not be validated"
}

func (e ErrValidation) Is(target error) bool {
	return target.Error() == "domain object could not be validated"
}

func (e ErrValidation) Convert() validator.ValidationErrors {
	var valiErrs validator.ValidationErrors
	if !errors.As(e.ValiErr, &valiErrs) {
		slog.Error("could not convert err to validationErrors", "error", e.ValiErr)
		return nil
	}

	return valiErrs
}
