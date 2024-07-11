package validation

import (
	"fmt"
	"log/slog"
	"reflect"
)

type Rule interface {
	IsViolated(val any) bool
	Violation() error
	ViolationForHumans(val string) error
}

func PasswordMatchConfirmRule(confirm string) PasswordMatchConfirm {
	return PasswordMatchConfirm{
		confirm,
	}
}

type PasswordMatchConfirm struct {
	confirm string
}

// ViolationForHumans implements Rule.
func (p PasswordMatchConfirm) ViolationForHumans(val string) error {
	return fmt.Errorf(
		"password and confirm password must match",
	)
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

// ViolationForHumans implements Rule.
func (r Required) ViolationForHumans(val string) error {
	return fmt.Errorf("%v needs to provided", val)
}

// IsViolated implements Rule.
func (r Required) IsViolated(v any) bool {
	return IsEmpty(v)
}

// Violation implements Rule.
func (r Required) Violation() error {
	return ErrIsRequired
}

var MinLenRule = func(required int) MinLen {
	return MinLen{
		required,
	}
}

var MustBeTrueRule MustBeTrue

type MustBeTrue bool

// IsViolated implements Rule.
func (m MustBeTrue) IsViolated(val any) bool {
	v, ok := val.(bool)
	if !ok {
		slog.Error("MustBeTruerule recevied a non boolean value")
		return true
	}

	return !v
}

// Violation implements Rule.
func (m MustBeTrue) Violation() error {
	return ErrMustBeTrue
}

// ViolationForHumans implements Rule.
func (m MustBeTrue) ViolationForHumans(val string) error {
	return fmt.Errorf(
		"%s needs to be 'true'",
		val,
	)
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

func (m MinLen) ViolationForHumans(val string) error {
	return fmt.Errorf(
		"%s needs to be longer than: '%v' characters",
		val,
		m.requiredLen,
	)
}

var MaxLenRule = func(required int) MaxLen {
	return MaxLen{required}
}

type MaxLen struct {
	requiredLen int
}

// ViolationForHumans implements Rule.
func (m MaxLen) ViolationForHumans(val string) error {
	return fmt.Errorf(
		"%s needs to be shorter than: '%v' characters",
		val,
		m.requiredLen,
	)
}

// IsViolated implements Rule.
func (m MaxLen) IsViolated(val any) bool {
	v := reflect.ValueOf(val)
	valLen, err := LengthOfValue(&v)
	if err != nil {
		slog.Error("could not get length of value for MaxLenRule", "error", err, "val", val)
		return true
	}

	return valLen > m.requiredLen
}

// Violation implements Rule.
func (m MaxLen) Violation() error {
	return fmt.Errorf(
		"failed error: '%v', because val was longer than: '%v' characters",
		ErrTooLong,
		m,
	)
}

var EmailRule Email

type Email struct{}

// ViolationForHumans implements Rule.
func (e Email) ViolationForHumans(val string) error {
	return fmt.Errorf("the provided email was not valid")
}

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
	_ Rule = new(PasswordMatchConfirm)
	_ Rule = new(Required)
	_ Rule = new(MinLen)
	_ Rule = new(MaxLen)
	_ Rule = new(Email)
	_ Rule = new(MustBeTrue)
)
