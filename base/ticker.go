package base

import "time"

func NewTicker(stop chan struct{}, f func()) {
	t := time.NewTimer(time.Second)
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
