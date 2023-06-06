package base

import (
	"fmt"
	"testing"
	"time"
)

func TestNewTicker(t *testing.T) {
	s := make(chan struct{})
	go NewTicker(s, func() {
		fmt.Println("xxxxxx")
	})

	time.Sleep(time.Second * 5)
}
