package domain

import (
	"fmt"
)

type KeyValueDomainError struct {
	message string
}

func (e KeyValueDomainError) Error() string {
	return e.message
}

func NewInvalidLengthError(value string, valueName string, minValue int, maxValue int) *KeyValueDomainError {
	return &KeyValueDomainError{
		fmt.Sprintf(
			"expected '%v' to have a value between %v and %v; got %v (len=%v)",
			valueName, minValue, maxValue, value, len(value)),
	}
}

type KeyNotFoundError struct {
	message string
}

func (e KeyNotFoundError) Error() string {
	return e.message
}

func NewKeyNotFoundError(key string) error {
	return &KeyNotFoundError{
		fmt.Sprintf("No value found with key '%v'", key),
	}
}

type KeyExistsError struct {
	message string
}

func (e KeyExistsError) Error() string {
	return e.message
}

func NewKeyExistsError(key string) error {
	return &KeyExistsError{
		fmt.Sprintf("Key '%v' already exists", key),
	}
}
