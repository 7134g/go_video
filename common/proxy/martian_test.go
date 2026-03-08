package proxy

import (
	"fmt"
	"testing"
)

func TestMartian(t *testing.T) {
	cfg := NewConfig()
	cfg.EnableMITM = true
	cfg.UpstreamProxy = "http://127.0.0.1:7890"
	cfg.TaskHandler = func(task VideoTask) error {
		fmt.Printf("Captured video: %s (%s)\n", task.URL, task.VideoType)
		return nil
	}

	p, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}

	if err := p.Start(); err != nil {
		t.Fatal(err)
	}
}
