package controller

import (
	"context"
	"encoding/json"
	"errors"
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

	repo *repository.TaskRepository
}

func GetController() *DownloadController {
	once.Do(func() {
		dirPwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		downloadController = &DownloadController{
			tasks:        make(map[uint]*DTask),
			pwd:          dirPwd,
			config:       model.DefaultConfig(),
			taskQueue:    make(chan *DTask, 100),
			downloadPool: downloader.NewPool(model.DefaultConfig().MaxSegmentWorkers),
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
	delete(c.tasks, id)
	return nil
}

func (c *DownloadController) PauseTask(id uint) error {
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
	}
	task.cancel()
	return nil
}

type TaskCallback func(id uint, err error) error

func (c *DownloadController) StartAll(callback TaskCallback) {
	c.mu.RLock()
	tasks := make([]*DTask, 0, len(c.tasks))
	for _, t := range c.tasks {
		tasks = append(tasks, t)
	}
	maxConcurrent := c.config.MaxConcurrentTasks
	c.mu.RUnlock()

	sem := make(chan struct{}, maxConcurrent)
	for _, task := range tasks {
		sem <- struct{}{}
		go func(t *DTask) {
			defer func() { <-sem }()
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
	go c.runTask(task, callback)
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
