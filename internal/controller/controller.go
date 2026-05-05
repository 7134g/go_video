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

type DownloadController struct {
	mu           sync.RWMutex
	tasks        map[uint]*DTask
	pwd          string
	config       *model.Config
	runningCount int
	taskQueue    chan *DTask
	downloadPool *downloader.Pool
	taskSem      chan struct{}

	repo *repository.TaskRepository
}

func GetController() *DownloadController {
	once.Do(func() {
		dirPwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		cfg := model.DefaultConfig()
		downloadController = &DownloadController{
			tasks:        make(map[uint]*DTask),
			pwd:          dirPwd,
			config:       cfg,
			taskQueue:    make(chan *DTask, 100),
			downloadPool: downloader.NewPool(cfg.MaxSegmentWorkers),
			taskSem:      make(chan struct{}, cfg.MaxConcurrentTasks),
			repo:         repository.NewTaskRepository(),
		}
		if err := downloadController.repo.ResetStatus(); err != nil {
			panic(err)
		}
	})
	return downloadController
}

func (c *DownloadController) ReloadConfig() {
	// Called by service layer after config update
}

func (c *DownloadController) ApplyConfig(downloadDir string, maxConcurrent, maxSegment, maxErrors int, defaultHeaders map[string]string) {
	c.mu.Lock()
	c.config.DownloadDir = downloadDir
	c.config.MaxConcurrentTasks = maxConcurrent
	c.config.MaxSegmentWorkers = maxSegment
	c.config.MaxConsecutiveErrors = maxErrors
	c.config.DefaultHeaders = defaultHeaders
	c.downloadPool = downloader.NewPool(maxSegment)
	c.taskSem = make(chan struct{}, maxConcurrent)
	c.mu.Unlock()
}

func (c *DownloadController) SetDownloadDir(dir string) {
	c.config.DownloadDir = dir
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
		Progress: &Progress{},
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
		Progress: &Progress{},
		ctx:      ctx,
		cancel:   cancel,
	}

	c.mu.Lock()
	c.tasks[id] = task
	c.mu.Unlock()

	go c.runTask(task, callback)
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
		go func(t *DTask) {
			c.taskSem <- struct{}{}
			defer func() { <-c.taskSem }()
			c.runTask(t, callback)
		}(task)
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
	go func() {
		c.taskSem <- struct{}{}
		defer func() { <-c.taskSem }()
		c.runTask(task, callback)
	}()
	return nil
}

func (c *DownloadController) runTask(task *DTask, callback TaskCallback) {
	var err error
	switch task.Type {
	case TaskTypeMp4:
		err = c.downloadMp4(task)
	case TaskTypeM3u8:
		err = c.downloadM3u8(task)
		if err == nil {
			err = c.mergeM3u8(task)
			if err != nil {
				BroadcastMessage(task.ID, "合并失败..."+err.Error())
			}
		}
	}
	if callback != nil {
		callback(task.ID, err)
	}
}

type ProgressInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Downloaded  int64  `json:"downloaded"`
	Total       int64  `json:"total"`
	SegmentDone int    `json:"segment_done"`
	SegmentAll  int    `json:"segment_all"`
	Percent     int    `json:"percent"`
}

func (c *DownloadController) GetAllProgress() []ProgressInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make([]ProgressInfo, 0, len(c.tasks))
	for _, t := range c.tasks {
		downloaded, total := t.Progress.GetProgress()
		segDone, segAll := t.Progress.GetSegment()
		percent := 0
		if total > 0 {
			percent = int(downloaded * 100 / total)
		} else if segAll > 0 {
			percent = segDone * 100 / segAll
		}
		result = append(result, ProgressInfo{
			ID:          t.ID,
			Name:        t.Name,
			Type:        string(t.Type),
			Downloaded:  downloaded,
			Total:       total,
			SegmentDone: segDone,
			SegmentAll:  segAll,
			Percent:     percent,
		})
	}
	return result
}
