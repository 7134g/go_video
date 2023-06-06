package base

import (
	"time"
)

func NewTicker(stop chan struct{}, f func()) {
	t := time.NewTicker(time.Second * 3)
	defer t.Stop()
	for {
		select {
		case <-stop:
			return
		case <-t.C:
			f()
		}
	}
}
