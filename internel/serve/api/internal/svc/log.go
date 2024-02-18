package svc

import "bytes"

type logCache struct {
	Cache chan []byte
	bytes.Buffer
}

func newLogCache() *logCache {
	return &logCache{
		Cache: make(chan []byte, 200),
	}
}

func (l *logCache) Write(p []byte) (n int, err error) {
	if len(l.Cache) > 100 {
		<-l.Cache
	}
	l.Cache <- p
	return l.Buffer.Write(p)
}
