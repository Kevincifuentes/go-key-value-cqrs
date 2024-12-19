package persistence

import "sync"

type KeyValueMap struct {
	mu sync.RWMutex
	m  map[string]string
}

func NewKeyValueMap() *KeyValueMap {
	return &KeyValueMap{}
}

func (m *KeyValueMap) Find(key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok := m.m[key]
	return value
}
