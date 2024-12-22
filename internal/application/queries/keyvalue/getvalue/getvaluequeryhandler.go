package getvalue

import "go-key-value-cqrs/domain"

type QueryHandler struct {
	KeyValueReader domain.KeyValueReader
}

func (handler QueryHandler) Ask(query Query) (domain.KeyValueView, error) {
	return handler.KeyValueReader.Get(query.Key)
}
