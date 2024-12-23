package model

import "fmt"

type ApiValidationError struct {
	message string
}

func (e ApiValidationError) Error() string {
	return e.message
}

func NewApiValidationError(request interface{}, reason string) ApiValidationError {
	return ApiValidationError{
		fmt.Sprintf(
			"Invalid request '%v'. Reason: %v", request, reason),
	}

}
