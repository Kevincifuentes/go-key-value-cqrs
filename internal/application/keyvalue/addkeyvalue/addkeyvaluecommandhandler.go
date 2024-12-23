package addkeyvalue

import "go-key-value-cqrs/domain"

type CommandHandler struct {
	KeyValueWriter domain.KeyValueWriter
}

func (handler CommandHandler) Execute(command Command) error {
	newKeyValue, domainError := domain.NewKeyValue(command.Key, command.Value)
	if domainError != nil {
		return domainError
	}
	return handler.KeyValueWriter.Add(*newKeyValue)
}
