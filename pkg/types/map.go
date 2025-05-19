package types

import (
	"sync"
)

type SyncMap[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{
		mu:   sync.RWMutex{},
		data: make(map[K]V),
	}
}

func (sm *SyncMap[K, V]) Set(key K, value V) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.data[key] = value
}

func (sm *SyncMap[K, V]) Get(key K) (V, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	v, ok := sm.data[key]
	return v, ok
}

func (sm *SyncMap[K, V]) Pull(key K) (V, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	v, ok := sm.data[key]
	if ok {
		delete(sm.data, key)
	}
	return v, ok
}

func (sm *SyncMap[K, V]) Delete(key K) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	delete(sm.data, key)
}

func (sm *SyncMap[K, V]) Clear() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	clear(sm.data)
}
