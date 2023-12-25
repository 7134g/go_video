package table

import (
	"sync"
)

type dataType interface {
	~uint
}

type base[D dataType] struct {
	lock sync.RWMutex

	body map[string]D
}

func (m *base[D]) Set(key string, value D) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.body[key] = value
}

func (m *base[D]) Get(key string) (D, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	value, exist := m.body[key]
	return value, exist
}
