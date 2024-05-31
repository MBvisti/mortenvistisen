package usecases

import "errors"

type ValidationErrorsMap = map[string]string

var ErrValidationErrorConversion = errors.New("could not convert error to ValidationErrors")
