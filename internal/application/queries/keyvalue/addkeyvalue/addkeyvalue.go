package addkeyvalue

import (
	"go-key-value-cqrs/application/queries/cqrs/commandbus"
)

type Command struct {
	Key   string
	Value string
}

func (command Command) Config() commandbus.CommandConfig {
	return commandbus.CommandConfig{Name: "AddKeyValueCommand"}
}
