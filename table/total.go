package table

import "sync"

var M3U8DownloadSpeed = Speed{
	lock: sync.RWMutex{},
	body: make(map[string]uint),
}

type Speed struct {
	lock sync.RWMutex

	body map[string]uint
}

func (s *Speed) Increase(key string, value uint) {
	s.lock.Lock()
	defer s.lock.Unlock()

	v, ok := s.body[key]
	if ok {
		s.body[key] = v + value
	} else {
		s.body[key] = value
	}

}

func (s *Speed) Get(key string) (uint, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	v, ok := s.body[key]
	return v, ok
}
