package getvalue

import (
	"go-key-value-cqrs/application/cqrs/querybus"
)

type Query struct {
	Key string
}

func (query Query) Config() querybus.QueryConfig {
	return querybus.QueryConfig{Name: "GetKeyValueQuery"}
}
