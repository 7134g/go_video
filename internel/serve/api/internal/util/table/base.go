package table

import (
	"sync"
)

type dataType interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 | string
}

type cmpMap[D dataType] struct {
	lock sync.RWMutex

	body map[string]D
}

func (m *cmpMap[D]) Set(key string, value D) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.body[key] = value
}

func (m *cmpMap[D]) Get(key string) (D, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	value, exist := m.body[key]
	return value, exist
}

func (m *cmpMap[D]) Inc(key string, count D) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	value, exist := m.body[key]
	if exist {
		m.body[key] = count + value
	} else {
		m.body[key] = count
	}
}

type sliceType interface {
	[]byte | []string | []int
}

type sliceMap[D sliceType] struct {
	lock sync.RWMutex

	body map[string]D
}

func (m *sliceMap[D]) Set(key string, value D) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.body[key] = value
}

func (m *sliceMap[D]) Get(key string) (D, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	value, exist := m.body[key]
	return value, exist
}
