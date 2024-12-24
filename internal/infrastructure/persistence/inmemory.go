package persistence

import (
	"go-key-value-cqrs/domain"
	"sync"
)

// KeyValueMap Using Mutex because I'm considering frequent writes
type KeyValueMap struct {
	keyToValueMap map[string]string
	barrier       sync.RWMutex
}

type InMemoryKeyValueRepository struct {
	KeyValueMap KeyValueMap
}

func NewInMemoryKeyValueRepository() *InMemoryKeyValueRepository {
	return &InMemoryKeyValueRepository{
		KeyValueMap: KeyValueMap{
			keyToValueMap: make(map[string]string),
		},
	}
}

func (repository *InMemoryKeyValueRepository) Get(key string) (domain.KeyValueView, error) {
	barrier := &repository.KeyValueMap.barrier
	barrier.RLock()
	defer barrier.RUnlock()

	value, ok := repository.KeyValueMap.keyToValueMap[key]
	if !ok {
		return domain.KeyValueView{}, domain.NewKeyNotFoundError(key)
	}
	return domain.KeyValueView{Key: key, Value: value}, nil
}

func (repository *InMemoryKeyValueRepository) Add(keyValue domain.KeyValue) error {
	barrier := &repository.KeyValueMap.barrier
	barrier.Lock()
	defer barrier.Unlock()

	keyToAdd := keyValue.Key.Key
	if _, exists := repository.KeyValueMap.keyToValueMap[keyToAdd]; exists {
		return domain.NewKeyExistsError(keyToAdd)
	}
	repository.KeyValueMap.keyToValueMap[keyToAdd] = keyValue.Value.Value
	return nil
}

func (repository *InMemoryKeyValueRepository) Delete(key string) error {
	barrier := &repository.KeyValueMap.barrier
	barrier.Lock()
	defer barrier.Unlock()

	if _, exists := repository.KeyValueMap.keyToValueMap[key]; !exists {
		return domain.NewKeyNotFoundError(key)
	}

	delete(repository.KeyValueMap.keyToValueMap, key)
	return nil
}
