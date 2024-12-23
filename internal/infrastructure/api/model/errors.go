package model

import "fmt"

type ApiValidationError struct {
	message string
}

func (e ApiValidationError) Error() string {
	return e.message
}

type InvalidRequestError struct {
	ApiValidationError
}

func (e InvalidRequestError) Error() string {
	return e.message
}

func NewInvalidRequestError(request interface{}, reason string) error {
	return &InvalidRequestError{
		ApiValidationError{
			fmt.Sprintf(
				"Invalid request '%v'. Reason: %v", request, reason),
		},
	}
}
