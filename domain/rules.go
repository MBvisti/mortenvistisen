package domain

import (
	"fmt"
	"log/slog"
	"reflect"
)

type Rule interface {
	IsViolated(val any) bool
	Violation() error
}

// type Compareable interface {
// 	NotEqual(val, candidate any) bool
// }

func PasswordMatchConfirmRule(confirm string) PasswordMatchConfirm {
	return PasswordMatchConfirm{
		confirm,
	}
}

type PasswordMatchConfirm struct {
	confirm string
}

// Violation implements Rule.
func (p PasswordMatchConfirm) Violation() error {
	return ErrPasswordDontMatch
}

// IsViolated implements Rule.
func (p PasswordMatchConfirm) IsViolated(val any) bool {
	valString := fmt.Sprintf("%v", val)

	return valString != p.confirm
}

var RequiredRule Required

type Required struct{}

// IsViolated implements Rule.
func (r Required) IsViolated(v any) bool {
	// if v.Type().String() == "uuid.UUID" {
	// 	id, ok := v.(uuid.UUID)
	// 	if !ok {
	// 		slog.Error("could not convert val to uuid in Required")
	// 		return true
	// 	}
	//
	// 	return id == uuid.Nil
	// }
	return IsEmpty(v)
}

// Violation implements Rule.
func (r Required) Violation() error {
	return ErrIsRequired
}

var MinLenRule = func(required int) MinLen {
	return MinLen{
		requiredLen: required,
	}
}

type MinLen struct {
	requiredLen int
}

// IsViolated implements Rule.
func (m MinLen) IsViolated(val any) bool {
	v := reflect.ValueOf(val)
	valLen, err := LengthOfValue(&v)
	if err != nil {
		slog.Error("could not get length of value for MinLenRule", "error", err, "val", val)
		return true
	}

	return valLen < m.requiredLen
}

// Violation implements Rule.
func (m MinLen) Violation() error {
	return ErrTooShort
}

type MaxLenRule int

// IsViolated implements Rule.
func (m MaxLenRule) IsViolated(val any) bool {
	v := reflect.ValueOf(val)
	valLen, err := LengthOfValue(&v)
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

// _ Rule = new(PasswordMatchConfirm)
// _ Compareable = new(PasswordMatchConfirm)
// _ Rule = new(Required)
// var _ Rule = new(MinLenRule)

// _ Rule = new(MaxLenRule)
// _ Rule = new(Email)
