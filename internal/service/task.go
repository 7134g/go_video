package service

import (
	"context"
	"encoding/json"
	"errors"

	"go_video/internal/controller"
	"go_video/internal/model"
	"go_video/internal/repository"
)

type TaskService struct {
	repo *repository.TaskRepository
	ctrl *controller.DownloadController
}

func NewTaskService() *TaskService {
	return &TaskService{
		repo: repository.NewTaskRepository(),
		ctrl: controller.GetController(),
	}
}

func (s *TaskService) Create(task *model.Task) error {
	return s.repo.Create(task)
}

func (s *TaskService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *TaskService) Update(task *model.Task) error {
	return s.repo.Update(task)
}

func (s *TaskService) GetByID(id uint) (*model.Task, error) {
	return s.repo.GetByID(id)
}

func (s *TaskService) GetAll() ([]model.Task, error) {
	return s.repo.GetAll()
}

func (s *TaskService) GetByStatus(status model.TaskStatus) ([]model.Task, error) {
	return s.repo.GetByStatus(status)
}

func (s *TaskService) StartTasks() (int, error) {
	tasks, err := s.repo.GetPending()
	if err != nil {
		return 0, err
	}

	// Merge default headers from config into task headers
	cfg := GetConfigService().GetConfig()

	for _, t := range tasks {
		headerJSON := mergeHeaders(cfg.DefaultHeaders, t.Header)
		s.repo.UpdateStatus(t.ID, model.TaskStatusRunning)
		s.ctrl.AddTask(t.ID, t.Name, t.URL, headerJSON, t.Type)
	}

	s.ctrl.StartAll(func(id uint, err error) {
		if err != nil {
			if err == context.Canceled {
				s.repo.UpdateStatus(id, model.TaskStatusPaused)
			} else {
				s.repo.UpdateStatus(id, model.TaskStatusFailed)
			}
		} else {
			s.repo.UpdateStatus(id, model.TaskStatusCompleted)
		}
		s.ctrl.RemoveTask(id)
	})

	return len(tasks), nil
}

func (s *TaskService) PauseTask(id uint) error {
	if err := s.ctrl.PauseTask(id); err != nil {
		return err
	}
	return nil
}

func (s *TaskService) RetryTask(id uint) error {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if task.Status != model.TaskStatusFailed && task.Status != model.TaskStatusPaused {
		return errors.New("task is not in failed or paused state")
	}

	cfg := GetConfigService().GetConfig()
	headerJSON := mergeHeaders(cfg.DefaultHeaders, task.Header)

	s.repo.UpdateStatus(id, model.TaskStatusRunning)
	s.ctrl.AddTask(id, task.Name, task.URL, headerJSON, task.Type)

	s.ctrl.StartTask(id, func(tid uint, err error) {
		if err != nil {
			if err == context.Canceled {
				s.repo.UpdateStatus(tid, model.TaskStatusPaused)
			} else {
				s.repo.UpdateStatus(tid, model.TaskStatusFailed)
			}
		} else {
			s.repo.UpdateStatus(tid, model.TaskStatusCompleted)
		}
		s.ctrl.RemoveTask(tid)
	})

	return nil
}

// mergeHeaders merges default headers into task-specific headers.
// Task headers take priority over default headers.
func mergeHeaders(defaults map[string]string, taskHeaderJSON string) string {
	if len(defaults) == 0 {
		return taskHeaderJSON
	}

	merged := make(map[string][]string)
	for k, v := range defaults {
		merged[k] = []string{v}
	}

	if taskHeaderJSON != "" {
		var taskHeaders map[string][]string
		if err := json.Unmarshal([]byte(taskHeaderJSON), &taskHeaders); err == nil {
			for k, v := range taskHeaders {
				merged[k] = v
			}
		}
	}

	data, _ := json.Marshal(merged)
	return string(data)
}
