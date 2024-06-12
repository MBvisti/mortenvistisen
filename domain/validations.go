package domain

import (
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
)

const (
	UUIDType   = "uuid.UUID"
	StringType = "string"
	StructType = "struct"
	TimeType   = "time.Time"
	NilType    = "nil"
)

func LengthOfValue(v *reflect.Value) (int, error) {
	switch v.Kind() {
	case reflect.String, reflect.Slice, reflect.Map, reflect.Array:
		return v.Len(), nil
	case reflect.Invalid:
		return 0, fmt.Errorf("cannot get the length of invalid kind")
	default:
		return 0, fmt.Errorf("provided value: '%v' did not match any kind", v.Kind())
	}
}

func GetTypeOfValue(rfVal reflect.Value) any {
	var conv any
	switch rfVal.Type().String() {
	case StringType:
		convVal, _ := rfVal.Interface().(string)
		conv = convVal
	case UUIDType:
		convVal, _ := rfVal.Interface().(uuid.UUID)
		conv = convVal
	case StructType:
		convVal, _ := rfVal.Interface().(time.Time)
		conv = convVal
	case TimeType:
		convVal, _ := rfVal.Interface().(time.Time)
		conv = convVal
	case NilType:
	}

	return conv
}

func IsEmpty(value interface{}) bool {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Map, reflect.Slice:
		if v.Type().String() == "uuid.UUID" {
			uid, ok := v.Interface().(uuid.UUID)
			if !ok {
				return true
			}

			return uid == uuid.UUID{}
		}
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

func checkRule(val *reflect.Value, rule Rule) error {
	if isViolated := rule.IsViolated(val); isViolated {
		return rule.Violation()
	}

	return nil
}

func GetFieldValue(fieldValue reflect.Value) any {
	if !fieldValue.CanInterface() {
		return nil
	}
	return fieldValue.Interface()
}
