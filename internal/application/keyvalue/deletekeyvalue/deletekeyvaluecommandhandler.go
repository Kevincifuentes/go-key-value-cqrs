package deletekeyvalue

import "go-key-value-cqrs/domain"

type CommandHandler struct {
	KeyValueWriter domain.KeyValueWriter
}

func (handler CommandHandler) Execute(command Command) error {
	return handler.KeyValueWriter.Delete(command.Key)
}
