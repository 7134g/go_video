package proxy

import "time"

type VideoTask struct {
	URL     string
	Method  string
	Headers string
	Body    []byte
	Title   string
	Type    string

	CreateAt time.Time
}
