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
		if err := s.repo.UpdateStatus(t.ID, model.TaskStatusRunning); err != nil {
			return 0, err
		}
		if err := s.ctrl.AddTask(t.ID, t.Name, t.URL, headerJSON, t.Type); err != nil {
			return 0, err
		}
	}

	s.ctrl.StartAll(func(id uint, err error) error {
		defer func(ctrl *controller.DownloadController, id uint) {
			_ = ctrl.RemoveTask(id)
		}(s.ctrl, id)

		if err != nil {
			if errors.Is(err, context.Canceled) {
				return s.repo.UpdateStatus(id, model.TaskStatusPaused)
			} else {
				return s.repo.UpdateStatus(id, model.TaskStatusFailed)
			}
		} else {
			return s.repo.UpdateStatus(id, model.TaskStatusCompleted)
		}
	})

	return len(tasks), nil
}

func (s *TaskService) PauseTask(id uint) error {
	if err := s.ctrl.StopTask(id); err != nil {
		return err
	}
	return nil
}

func (s *TaskService) PauseAllTasks() error {
	s.ctrl.StopAll()
	return s.repo.UpdateAllByStatus(model.TaskStatusRunning, model.TaskStatusPaused)
}

func (s *TaskService) StartTask(id uint) error {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if task.Status != model.TaskStatusPending && task.Status != model.TaskStatusFailed && task.Status != model.TaskStatusPaused {
		return errors.New("task is already running or completed")
	}

	cfg := GetConfigService().GetConfig()
	headerJSON := mergeHeaders(cfg.DefaultHeaders, task.Header)

	if err := s.repo.UpdateStatus(id, model.TaskStatusRunning); err != nil {
		return err
	}
	if err := s.ctrl.AddTask(id, task.Name, task.URL, headerJSON, task.Type); err != nil {
		return err
	}

	return s.ctrl.StartTask(id, func(id uint, err error) error {
		defer func(ctrl *controller.DownloadController, id uint) {
			_ = ctrl.RemoveTask(id)
		}(s.ctrl, id)

		if err != nil {
			if errors.Is(err, context.Canceled) {
				return s.repo.UpdateStatus(id, model.TaskStatusPaused)
			}
			return s.repo.UpdateStatus(id, model.TaskStatusFailed)
		}
		return s.repo.UpdateStatus(id, model.TaskStatusCompleted)
	})
}

func (s *TaskService) AddAndStart(task *model.Task) error {
	if err := s.repo.Create(task); err != nil {
		return err
	}

	cfg := GetConfigService().GetConfig()
	headerJSON := mergeHeaders(cfg.DefaultHeaders, task.Header)

	if err := s.repo.UpdateStatus(task.ID, model.TaskStatusRunning); err != nil {
		return err
	}

	return s.ctrl.AddAndStart(task.ID, task.Name, task.URL, headerJSON, task.Type, func(id uint, err error) error {
		defer func(ctrl *controller.DownloadController, id uint) {
			_ = ctrl.RemoveTask(id)
		}(s.ctrl, id)

		if err != nil {
			if errors.Is(err, context.Canceled) {
				return s.repo.UpdateStatus(id, model.TaskStatusPaused)
			}
			return s.repo.UpdateStatus(id, model.TaskStatusFailed)
		}
		return s.repo.UpdateStatus(id, model.TaskStatusCompleted)
	})
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

	if err := s.repo.UpdateStatus(id, model.TaskStatusRunning); err != nil {
		return err
	}
	if err := s.ctrl.AddTask(id, task.Name, task.URL, headerJSON, task.Type); err != nil {
		return err
	}

	s.ctrl.StartAll(func(id uint, err error) error {
		defer func(ctrl *controller.DownloadController, id uint) {
			_ = ctrl.RemoveTask(id)
		}(s.ctrl, id)

		if err != nil {
			if errors.Is(err, context.Canceled) {
				return s.repo.UpdateStatus(id, model.TaskStatusPaused)
			} else {
				return s.repo.UpdateStatus(id, model.TaskStatusFailed)
			}
		} else {
			return s.repo.UpdateStatus(id, model.TaskStatusCompleted)
		}
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
