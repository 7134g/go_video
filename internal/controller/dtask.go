package controller

import (
	"context"
	"net/http"
	"sync"
)

type TaskType string

const (
	TaskTypeMp4  TaskType = "mp4"
	TaskTypeM3u8 TaskType = "m3u8"
)

type DTask struct {
	ID       uint
	Name     string
	URL      string
	Header   http.Header
	Type     TaskType
	Progress *Progress
	ctx      context.Context
	cancel   context.CancelFunc
}

type Progress struct {
	mu          sync.RWMutex
	Downloaded  int64
	Total       int64
	SegmentDone int
	SegmentAll  int
	Status      string
}

func (p *Progress) SetTotal(total int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Total = total
}

func (p *Progress) AddDownloaded(n int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Downloaded += n
}

func (p *Progress) GetProgress() (downloaded, total int64) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Downloaded, p.Total
}

func (p *Progress) SetSegment(done, all int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.SegmentDone = done
	p.SegmentAll = all
}

func (p *Progress) GetSegment() (done, all int) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.SegmentDone, p.SegmentAll
}
