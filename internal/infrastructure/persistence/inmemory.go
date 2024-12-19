package persistence

import (
	"go-key-value-cqrs/domain"
	"sync"
)

type KeyValueMap struct {
	keyToValueMap map[string]string
	barrier       sync.RWMutex
}

type InMemoryKeyValueRepository struct {
	keyValueMap KeyValueMap
}

func (repository *InMemoryKeyValueRepository) Get(key string) (domain.KeyValueView, error) {
	barrier := &repository.keyValueMap.barrier
	barrier.RLock()
	defer barrier.RUnlock()
	value, ok := repository.keyValueMap.keyToValueMap[key]
	if !ok {
		return domain.KeyValueView{}, domain.NewKeyNotFoundError(key)
	}
	return domain.KeyValueView{Key: key, Value: value}, nil
}
