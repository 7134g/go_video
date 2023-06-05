package table

import "sync"

var CryptoVedioTable CryptoTable

type CryptoTable struct {
	lock sync.RWMutex

	body map[string][]byte
}

func (c *CryptoTable) Set(key string, value []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.body[key] = value
}

func (c *CryptoTable) Get(key string) []byte {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.body[key]
}
