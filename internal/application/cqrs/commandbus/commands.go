package commandbus

type Command interface {
	Config() CommandConfig
}

type CommandConfig struct {
	Name string
}
