package domain

import (
	"fmt"
)

type KeyValueError struct {
	message string
}

type InvalidLengthError struct {
	KeyValueError
}

func (e InvalidLengthError) Error() string {
	return e.message
}

func NewInvalidLengthError(value string, valueName string, minValue int, maxValue int) error {
	return &InvalidLengthError{
		KeyValueError{
			fmt.Sprintf(
				"expected '%v' to have a value between %v and %v; got %v (len=%v)",
				valueName, minValue, maxValue, value, len(value)),
		},
	}
}

type KeyNotFoundError struct {
	KeyValueError
}

func (e KeyNotFoundError) Error() string {
	return e.message
}

func NewKeyNotFoundError(key string) error {
	return &KeyNotFoundError{
		KeyValueError{
			fmt.Sprintf("No value found with key '%v'", key),
		},
	}
}
