package querybus

import (
	"reflect"
	"sync"
)

// Using sync.Map because we will initial few writes and after that only reads
var registeredQueries = &sync.Map{}

func Load[T Query, K any](handler QueryHandler[T, K]) error {
	var query T
	registeredQueries.Store(query.Config().Name, handler)
	return nil
}

func Asks[K any, T Query](query T) (value K, err error) {
	queryName := query.Config().Name
	handler, ok := registeredQueries.Load(queryName)
	if !ok {
		return value, NewErrorQueryNotFound(queryName)
	}
	typedHandler, ok := handler.(QueryHandler[T, K])
	if !ok {
		return value, NewErrorQueryHandlerTypeNotValid(reflect.TypeOf(handler).Name(), queryName)
	}
	return typedHandler.Ask(query)
}
