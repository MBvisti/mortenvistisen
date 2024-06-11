package domain

import (
	"log/slog"
	"reflect"

	"github.com/google/uuid"
)

type Rule interface {
	IsViolated(val any) bool
	Violation() error
}

type Compareable interface {
	NotEqual(val, candidate any) bool
}

var PasswordMatchConfirmRule PasswordMatchConfirm

type PasswordMatchConfirm struct{}

// Compare implements Compareable.
func (p PasswordMatchConfirm) NotEqual(val any, candidate any) bool {
	valString, err := ToString(val)
	if err != nil {
		return true
	}

	candidateString, err := ToString(candidate)
	if err != nil {
		return true
	}

	return valString != candidateString
}

// Violation implements Rule.
func (p PasswordMatchConfirm) Violation() error {
	return ErrPasswordDontMatch
}

// IsViolated implements Rule.
func (p PasswordMatchConfirm) IsViolated(val any) bool {
	panic("unimplemented")
}

var RequiredRule Required

type Required struct{}

// IsViolated implements Rule.
func (r Required) IsViolated(val any) bool {
	if v := reflect.ValueOf(val); v.Type().String() == "uuid.UUID" {
		id, ok := val.(uuid.UUID)
		if !ok {
			slog.Error("could not convert val to uuid in Required")
			return true
		}

		return id == uuid.Nil
	}
	return IsEmpty(val)
}

// Violation implements Rule.
func (r Required) Violation() error {
	return ErrIsRequired
}

type MinLenRule int

// IsViolated implements Rule.
func (m MinLenRule) IsViolated(val any) bool {
	valLen, err := LengthOfValue(val)
	if err != nil {
		slog.Error("could not get length of value for MinLenRule", "error", err, "val", val)
		return true
	}

	return valLen < int(m)
}

// Violation implements Rule.
func (m MinLenRule) Violation() error {
	return ErrTooShort
}

type MaxLenRule int

// IsViolated implements Rule.
func (m MaxLenRule) IsViolated(val any) bool {
	valLen, err := LengthOfValue(val)
	if err != nil {
		slog.Error("could not get length of value for MaxLenRule", "error", err, "val", val)
		return true
	}

	return valLen > int(m)
}

// Violation implements Rule.
func (m MaxLenRule) Violation() error {
	return ErrTooLong
}

var EmailRule Email

type Email struct{}

// IsViolated implements Rule.
func (e Email) IsViolated(val any) bool {
	stringVal, err := ToString(val)
	if err != nil {
		return true
	}

	return !isEmailValid(stringVal)
}

// Violation implements Rule.
func (e Email) Violation() error {
	return ErrInvalidEmail
}

var (
	_ Rule        = new(PasswordMatchConfirm)
	_ Compareable = new(PasswordMatchConfirm)
	_ Rule        = new(Required)
	_ Rule        = new(MinLenRule)
	_ Rule        = new(MaxLenRule)
	_ Rule        = new(Email)
)
