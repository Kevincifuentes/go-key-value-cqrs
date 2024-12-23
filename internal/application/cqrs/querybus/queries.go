package querybus

type Query interface {
	Config() QueryConfig
}

type QueryConfig struct {
	Name string
}
