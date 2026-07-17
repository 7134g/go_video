//go:build desktop || bindings

package main

import (
	"context"
	"go_video/pkg/downloader"

	"go_video/internal/controller"
	"go_video/internal/model"
	"go_video/internal/service"
	"go_video/pkg/proxy"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx    context.Context
	svc    *service.TaskService
	cfgSvc *service.ConfigService
	ctrl   *controller.DownloadController
}

func NewApp() *App {
	return &App{
		svc:    service.NewTaskService(),
		cfgSvc: service.GetConfigService(),
		ctrl:   controller.GetController(),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	go a.bridgeMessages()
	go a.bridgeProgress()
}

func (a *App) shutdown(ctx context.Context) {
	a.ctrl.StopAll()
	a.cfgSvc.Shutdown()
}

func (a *App) AddTask(name, url, header, taskType string) (*model.Task, error) {
	task := &model.Task{
		Name:   name,
		URL:    url,
		Header: header,
		Type:   taskType,
	}
	if err := a.svc.Create(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (a *App) GetTasks(status int) ([]model.Task, error) {
	if status >= 0 {
		return a.svc.GetByStatus(model.TaskStatus(status))
	}
	return a.svc.GetAll()
}

func (a *App) DeleteTask(id uint) error {
	return a.svc.Delete(id)
}

type UpdateTaskFields struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	URL    string `json:"url"`
	Header string `json:"header"`
	Type   string `json:"type"`
}

func (a *App) UpdateTask(id uint, name, url, header, taskType string) (*model.Task, error) {
	task, err := a.svc.GetByID(id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		task.Name = name
	}
	if url != "" {
		task.URL = url
	}
	task.Header = header
	if taskType != "" {
		task.Type = taskType
	}

	if err := a.svc.Update(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (a *App) StartOneTask(id uint) error {
	return a.svc.StartTask(id)
}

func (a *App) PauseTask(id uint) error {
	return a.svc.PauseTask(id)
}

func (a *App) RetryTask(id uint) error {
	return a.svc.RetryTask(id)
}

func (a *App) StopAllTasks() error {
	return a.svc.PauseAllTasks()
}

func (a *App) StartAllTasks() (int, error) {
	return a.svc.StartTasks()
}

func (a *App) GetConfig() *model.Config {
	return a.cfgSvc.GetConfig()
}

func (a *App) UpdateConfig(updates map[string]interface{}) (*model.Config, error) {
	return a.cfgSvc.UpdateConfig(updates)
}

func (a *App) GetFfmpegStatus() map[string]bool {
	return map[string]bool{
		"exists":    downloader.Exists(),
		"supported": downloader.Supported(),
	}
}

func (a *App) DownloadFfmpeg() map[string]bool {
	if err := downloader.Download(context.Background()); err != nil {
		return map[string]bool{"exists": false}
	}
	return map[string]bool{"exists": true}
}

func (a *App) UpdateTaskTitle(id uint) (*model.Task, error) {
	task, err := a.svc.GetByID(id)
	if err != nil {
		return nil, err
	}

	title := proxy.SearchTitle(task.URL)
	if title == "" {
		return task, nil
	}

	task.Name = title
	if err := a.svc.Update(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (a *App) RedownloadTask(id uint) error {
	return a.svc.RedownloadTask(id)
}

func (a *App) CheckCaInstalled() map[string]any {
	installed, err := proxy.CheckCertInstalled()
	if err != nil {
		return map[string]any{"installed": false, "error": err.Error()}
	}
	return map[string]any{"installed": installed}
}

func (a *App) GetAllProgress() []controller.ProgressInfo {
	return a.ctrl.GetAllProgress()
}

func (a *App) bridgeMessages() {
	ch := make(chan controller.Message, 10)
	controller.AddMessageListener(ch)
	defer controller.RemoveMessageListener(ch)

	for {
		select {
		case msg := <-ch:
			runtime.EventsEmit(a.ctx, "task:broadcast", msg)
		case <-a.ctx.Done():
			return
		}
	}
}

func (a *App) bridgeProgress() {
	ch := make(chan controller.ProgressInfo, 10)
	controller.AddProgressListener(ch)
	defer controller.RemoveProgressListener(ch)

	for {
		select {
		case info := <-ch:
			runtime.EventsEmit(a.ctx, "task:progress", info)
		case <-a.ctx.Done():
			return
		}
	}
}
