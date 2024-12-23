package commandbus

import (
	"reflect"
	"sync"
)

// Using sync.Map because we will initial few writes and after that only reads
var registeredCommands = &sync.Map{}

func Load[T Command](handler CommandHandler[T]) error {
	var command T
	registeredCommands.Store(command.Config().Name, handler)
	return nil
}

func Execute[T Command](command T) (err error) {
	commandName := command.Config().Name
	handler, ok := registeredCommands.Load(commandName)
	if !ok {
		return NewErrorCommandNotFound(commandName)
	}
	typedHandler, ok := handler.(CommandHandler[T])
	if !ok {
		return NewErrorCommandHandlerTypeNotValid(reflect.TypeOf(handler).Name(), commandName)
	}
	return typedHandler.Execute(command)
}
