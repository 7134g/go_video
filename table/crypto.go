package table

import "sync"

var CryptoVideoTable = CryptoTable{
	lock: sync.RWMutex{},
	body: make(map[string][]byte),
}

type CryptoTable struct {
	lock sync.RWMutex

	body map[string][]byte
}

func (c *CryptoTable) Set(key string, value []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.body[key] = value
}

func (c *CryptoTable) Get(key string) ([]byte, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	v, ok := c.body[key]
	return v, ok
}
