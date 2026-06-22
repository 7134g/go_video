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

// DTask 是 controller 持有的运行时任务表示，区别于持久化的 model.Task。
// ctx/cancel 用于响应 StopTask / StopAll；callback 在终态时由 dispatch 调用。
type DTask struct {
	ID       uint
	Name     string
	URL      string
	Header   http.Header
	Type     TaskType
	Progress *Progress
	callback TaskCallback
	ctx      context.Context
	cancel   context.CancelFunc
}

// Progress 在多 segment goroutine 并发更新下保证安全；所有读写均通过其方法走 mu。
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
	if done > p.SegmentDone {
		p.SegmentDone = done
	}
	p.SegmentAll = all
}

func (p *Progress) IncrementSegmentDone() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.SegmentDone++
}

func (p *Progress) GetSegment() (done, all int) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.SegmentDone, p.SegmentAll
}
