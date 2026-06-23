package controller

import (
	"context"
	"net/http"
	"sync"
	"time"
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
	mu       sync.RWMutex
	taskID   uint
	taskName string
	Done     int64
	Total    int64
	Type     TaskType
	Status   string
}

func (p *Progress) SetTotal(total int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Total = total
}

func (p *Progress) AddDone(n int64) {
	p.mu.Lock()
	p.Done += n
	done, total := p.Done, p.Total
	p.mu.Unlock()

	percent := 0
	if total > 0 {
		percent = int(done * 100 / total)
	}
	BroadcastProgress(ProgressInfo{
		ID:       p.taskID,
		Name:     p.taskName,
		Type:     string(p.Type),
		Done:     done,
		Total:    total,
		Percent:  percent,
		Timespec: time.Now().UnixMilli(),
	})
}

func (p *Progress) IncrementDone() {
	p.mu.Lock()
	p.Done++
	done, total := p.Done, p.Total
	p.mu.Unlock()

	percent := 0
	if total > 0 {
		percent = int(done * 100 / total)
	}
	BroadcastProgress(ProgressInfo{
		ID:       p.taskID,
		Name:     p.taskName,
		Type:     string(p.Type),
		Done:     done,
		Total:    total,
		Percent:  percent,
		Timespec: time.Now().UnixMilli(),
	})
}

func (p *Progress) GetProgress() (done, total int64) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Done, p.Total
}
