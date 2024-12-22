package getvalue

import (
	"go-key-value-cqrs/application/queries/cqrs/querybus"
)

type Query struct {
	Key string
}

func (query Query) Config() querybus.QueryConfig {
	return querybus.QueryConfig{Name: "GetKeyValueQuery"}
}
