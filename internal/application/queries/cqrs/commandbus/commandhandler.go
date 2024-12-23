package commandbus

type CommandHandler[T Command] interface {
	Execute(command T) error
}
