package querybus

import "fmt"

type QueryError struct {
	message string
}

type ErrorQueryNotFound struct {
	QueryError
}

func (e ErrorQueryNotFound) Error() string {
	return e.message
}

func NewErrorQueryNotFound(queryName string) error {
	return &ErrorQueryNotFound{
		QueryError{
			fmt.Sprintf("No query found with name '%v'", queryName),
		},
	}
}

type ErrorQueryHandlerTypeNotValid struct {
	QueryError
}

func (e ErrorQueryHandlerTypeNotValid) Error() string {
	return e.message
}

func NewErrorQueryHandlerTypeNotValid(typeName string, queryName string) error {
	return &ErrorQueryHandlerTypeNotValid{
		QueryError{
			fmt.Sprintf("Handler has an incorrect type %s for Query %s", typeName, queryName),
		},
	}
}
