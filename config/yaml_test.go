package config

import (
	"fmt"
	"testing"
	"time"
)

type stu struct {
	count int
	Do    func()
}

func (s *stu) Name() {
	s.count++
}

func (s *stu) load() {
	s.Do = s.Name
}

func newStu() stu {
	s := stu{
		count: 0,
		Do:    nil,
	}
	s.Do = s.Name
	return s
}

func TestLoadYaml(t *testing.T) {
	s := newStu()
	go func() {
		fmt.Println(s)
		s.Do()
	}()
	time.Sleep(time.Second * 1)
	fmt.Println(s)
}
