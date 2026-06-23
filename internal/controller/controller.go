package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go_video/internal/downloader"
	"go_video/internal/model"
	"go_video/internal/repository"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	downloadController *DownloadController
	once               sync.Once
)

// DownloadController 是进程内单例的下载调度核心。
//
// 调度模型：
//   - 任务先进 taskQueue（缓冲 100），dispatch goroutine 顺序读取
//   - 每个任务通过 taskSem 信号量获取并发槽（容量 = MaxConcurrentTasks），再起独立 goroutine 跑
//   - 任务级取消通过 DTask.ctx；分段级并发由 downloadPool 按域名再次限流
//
// 注意：tasks map 只存"运行中/排队中"的任务；调用方完成后必须显式 RemoveTask。
type DownloadController struct {
	mu           sync.RWMutex
	tasks        map[uint]*DTask
	pwd          string
	config       *model.Config
	taskQueue    chan *DTask
	downloadPool *downloader.Pool
	taskSem      chan struct{}
	httpClient   *httpClientHolder

	repo *repository.TaskRepository
}

func GetController() *DownloadController {
	once.Do(func() {
		dirPwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		cfg := model.DefaultConfig()
		holder := &httpClientHolder{}
		holder.Init()
		if len(cfg.VpnAddress) > 0 && cfg.VpnStatus {
			holder.SetProxy(cfg.VpnAddress)
		}
		downloadController = &DownloadController{
			tasks:        make(map[uint]*DTask),
			pwd:          dirPwd,
			config:       cfg,
			taskQueue:    make(chan *DTask, 100),
			downloadPool: downloader.NewPool(cfg.MaxSegmentWorkers),
			taskSem:      make(chan struct{}, cfg.MaxConcurrentTasks),
			httpClient:   holder,
			repo:         repository.NewTaskRepository(),
		}
		if err := downloadController.repo.ResetStatus(); err != nil {
			panic(err)
		}
		go downloadController.dispatch()
	})
	return downloadController
}

func (c *DownloadController) ApplyConfig(downloadDir string, maxConcurrent, maxSegment, maxErrors int, defaultHeaders map[string]string) {
	c.mu.Lock()
	c.config.DownloadDir = downloadDir
	c.config.MaxConcurrentTasks = maxConcurrent
	c.config.MaxSegmentWorkers = maxSegment
	c.config.MaxConsecutiveErrors = maxErrors
	c.config.DefaultHeaders = defaultHeaders
	// NOTE: 替换 pool/taskSem 不影响已 in-flight 的任务（它们持有旧引用），
	// 新任务用新容量。完全一致需要更复杂的可调容量实现，目前权衡为简单。
	c.downloadPool = downloader.NewPool(maxSegment)
	c.taskSem = make(chan struct{}, maxConcurrent)
	c.mu.Unlock()
	if len(c.config.VpnAddress) > 0 && c.config.VpnStatus {
		c.httpClient.SetProxy(c.config.VpnAddress)
	}
}

func (c *DownloadController) SetDownloadDir(dir string) {
	c.config.DownloadDir = dir
}

// DeleteTaskFiles removes downloaded files for a completed task.
// For mp4: removes <name>.mp4; for m3u8: removes the whole segment directory.
// Returns nil if files don't exist.
func (c *DownloadController) DeleteTaskFiles(name string, taskType TaskType) error {
	base := safeJoin(c.config.DownloadDir, name)
	switch taskType {
	case TaskTypeMp4:
		_ = os.Remove(base + ".mp4")
	case TaskTypeM3u8:
		_ = os.RemoveAll(base)
	}
	return nil
}

func (c *DownloadController) AddTask(id uint, name, url, headerJSON, taskType string) error {
	header := http.Header{}
	if headerJSON != "" {
		if err := json.Unmarshal([]byte(headerJSON), &header); err != nil {
			return err
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	task := &DTask{
		ID:       id,
		Name:     name,
		URL:      url,
		Header:   header,
		Type:     TaskType(taskType),
		Progress: &Progress{Type: TaskType(taskType)},
		ctx:      ctx,
		cancel:   cancel,
	}

	c.mu.Lock()
	c.tasks[id] = task
	c.mu.Unlock()

	BroadcastMessage(id, "任务已添加: "+name)
	return nil
}

func (c *DownloadController) AddAndStart(id uint, name, url, headerJSON, taskType string, callback TaskCallback) error {
	header := http.Header{}
	if headerJSON != "" {
		if err := json.Unmarshal([]byte(headerJSON), &header); err != nil {
			return err
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	task := &DTask{
		ID:       id,
		Name:     name,
		URL:      url,
		Header:   header,
		Type:     TaskType(taskType),
		Progress: &Progress{Type: TaskType(taskType)},
		ctx:      ctx,
		cancel:   cancel,
	}

	c.mu.Lock()
	c.tasks[id] = task
	c.mu.Unlock()

	task.callback = callback
	c.taskQueue <- task
	BroadcastMessage(id, "任务已添加并启动: "+name)
	return nil
}

func (c *DownloadController) GetTask(id uint) *DTask {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.tasks[id]
}

func (c *DownloadController) RemoveTask(id uint) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	task, ok := c.tasks[id]
	if ok {
		BroadcastMessage(id, "任务已删除: "+task.Name)
	}
	delete(c.tasks, id)
	return nil
}

func (c *DownloadController) StopTask(id uint) error {
	c.mu.RLock()
	task, ok := c.tasks[id]
	c.mu.RUnlock()
	if !ok {
		task, err := c.repo.GetByID(id)
		if err != nil {
			log.Println(err)
			return err
		}
		if task == nil {
			return errors.New("task not found")
		}

		_ = c.repo.UpdateStatus(task.ID, model.TaskStatusPaused)
		return nil
	}
	name := task.Name
	task.cancel()
	BroadcastMessage(id, "任务已停止: "+name)
	return nil
}

func (c *DownloadController) StopAll() {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, t := range c.tasks {
		t.cancel()
	}
	BroadcastMessage(0, "已停止所有任务")
}

func (c *DownloadController) dispatch() {
	for task := range c.taskQueue {
		// 快照 taskSem，避免 ApplyConfig 替换后 release 到错误的 channel。
		c.mu.RLock()
		sem := c.taskSem
		c.mu.RUnlock()
		sem <- struct{}{}
		go func(t *DTask, sem chan struct{}) {
			defer func() { <-sem }()
			select {
			case <-t.ctx.Done():
				if t.callback != nil {
					_ = t.callback(t.ID, context.Canceled)
				}
			default:
				c.runTask(t, t.callback)
			}
		}(task, sem)
	}
}

// TaskCallback 在任务终态触达时被回调；返回 error 用于回写数据库状态时透传。
// 入参 err 可能是 context.Canceled（用户暂停）、io/网络错误（失败）或 nil（成功）。
type TaskCallback func(id uint, err error) error

func (c *DownloadController) StartAll(callback TaskCallback) {
	c.mu.RLock()
	tasks := make([]*DTask, 0, len(c.tasks))
	for _, t := range c.tasks {
		tasks = append(tasks, t)
	}
	c.mu.RUnlock()

	BroadcastMessage(0, "已启动 "+fmt.Sprint(len(tasks))+" 个任务")

	for _, task := range tasks {
		task.callback = callback
		c.taskQueue <- task
	}
}

func (c *DownloadController) StartTask(id uint, callback TaskCallback) error {
	c.mu.RLock()
	task, ok := c.tasks[id]
	c.mu.RUnlock()
	if !ok {
		return errors.New("task not found")
	}
	BroadcastMessage(id, "任务已启动: "+task.Name)
	task.callback = callback
	c.taskQueue <- task
	return nil
}

func (c *DownloadController) runTask(task *DTask, callback TaskCallback) {
	var err error
	switch task.Type {
	case TaskTypeMp4:
		err = c.downloadMp4(task)
		if err != nil {
			log.Println("下载失败---", err.Error())
			break
		}
	case TaskTypeM3u8:
		err = c.downloadM3u8(task)
		if err != nil {
			log.Println("下载失败...", err.Error())
			break
		}
		if err := c.mergeM3u8(task); err != nil {
			BroadcastMessage(task.ID, "合并失败..."+err.Error())
		}

	}
	if callback != nil {
		_ = callback(task.ID, err)
	}
}

type ProgressInfo struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Done    int64  `json:"done"`
	Total   int64  `json:"total"`
	Percent int    `json:"percent"`
}

func (c *DownloadController) GetAllProgress() []ProgressInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make([]ProgressInfo, 0, len(c.tasks))
	for _, t := range c.tasks {
		done, total := t.Progress.GetProgress()
		percent := 0
		if total > 0 {
			percent = int(done * 100 / total)
		}
		result = append(result, ProgressInfo{
			ID:      t.ID,
			Name:    t.Name,
			Type:    string(t.Type),
			Done:    done,
			Total:   total,
			Percent: percent,
		})
	}
	return result
}
