package commandbus

import "fmt"

type CommandError struct {
	message string
}

type ErrorCommandNotFound struct {
	CommandError
}

func (e ErrorCommandNotFound) Error() string {
	return e.message
}

func NewErrorCommandNotFound(commandName string) error {
	return &ErrorCommandNotFound{
		CommandError{
			fmt.Sprintf("No command found with name '%v'", commandName),
		},
	}
}

type ErrorCommandHandlerTypeNotValid struct {
	CommandError
}

func (e ErrorCommandHandlerTypeNotValid) Error() string {
	return e.message
}

func NewErrorCommandHandlerTypeNotValid(typeName string, commandName string) error {
	return &ErrorCommandHandlerTypeNotValid{
		CommandError{
			fmt.Sprintf("Handler has an incorrect type %s for Command %s", typeName, commandName),
		},
	}
}
