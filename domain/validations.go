package domain

import (
	"fmt"
	"reflect"
	"time"
)

const UUID = "uuid.UUID"

func LengthOfValue(value interface{}) (int, error) {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String, reflect.Slice, reflect.Map, reflect.Array:
		return v.Len(), nil
	}
	return 0, fmt.Errorf("cannot get the length of %v", v.Kind())
}

func IsEmpty(value interface{}) bool {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Map, reflect.Slice:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Invalid:
		return true
	case reflect.Interface, reflect.Ptr:
		if v.IsNil() {
			return true
		}
		return IsEmpty(v.Elem().Interface())
	case reflect.Struct:
		v, ok := value.(time.Time)
		if ok && v.IsZero() {
			return true
		}
	}

	return false
}

// ToFloat converts the given value to a float64.
// An error is returned for all incompatible types.
func ToFloat(value interface{}) (float64, error) {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		return v.Float(), nil
	}

	return 0, fmt.Errorf("cannot convert %v to float64", v.Kind())
}

// ToInt converts the given value to an int64.
// An error is returned for all incompatible types.
func ToInt(value interface{}) (int64, error) {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int(), nil
	}

	return 0, fmt.Errorf("cannot convert %v to int64", v.Kind())
}

// StringOrBytes typecasts a value into a string or byte slice.
// Boolean flags are returned to indicate if the typecasting succeeds or not.
func ToString(value interface{}) (string, error) {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		return v.String(), nil
	}

	return "", fmt.Errorf("cannot convert %v to string", v.Kind())
}

func checkRule(val any, rule Rule) error {
	if isViolated := rule.IsViolated(val); isViolated {
		return rule.Violation()
	}

	return nil
}

func checkComparableRule(
	val, compareVal any,
	rule Rule,
	eomparableRule Compareable,
) error {
	if isViolated := eomparableRule.NotEqual(val, compareVal); isViolated {
		return rule.Violation()
	}

	return nil
}
