package repository

import (
	"errors"

	"go_video/internal/model"

	"gorm.io/gorm"
)

type TaskRepository struct{}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{}
}

func (r *TaskRepository) Create(task *model.Task) error {
	return DB.Create(task).Error
}

func (r *TaskRepository) GetByID(id uint) (*model.Task, error) {
	var task model.Task
	if err := DB.First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) Update(task *model.Task) error {
	return DB.Save(task).Error
}

func (r *TaskRepository) Delete(id uint) error {
	return DB.Delete(&model.Task{}, id).Error
}

func (r *TaskRepository) GetPending() ([]model.Task, error) {
	var tasks []model.Task
	err := DB.Where("status IN ?", []model.TaskStatus{model.TaskStatusPending, model.TaskStatusPaused}).Find(&tasks).Error
	return tasks, err
}

// ResetStatus 用于进程启动时的崩溃恢复：上次进程被强杀时残留的 Running
// 实际上不再在跑，统一置回 Pending 等待用户重新触发。
func (r *TaskRepository) ResetStatus() error {
	return DB.Model(&model.Task{}).
		Where("status = ?", model.TaskStatusRunning).
		Update("status", model.TaskStatusPending).Error
}

func (r *TaskRepository) UpdateStatus(id uint, status model.TaskStatus) error {
	return DB.Model(&model.Task{}).Where("id = ?", id).Update("status", status).Error
}

func (r *TaskRepository) UpdateAllByStatus(from, to model.TaskStatus) error {
	return DB.Model(&model.Task{}).Where("status = ?", from).Update("status", to).Error
}

func (r *TaskRepository) GetAll() ([]model.Task, error) {
	var tasks []model.Task
	err := DB.Order("created_at desc").Find(&tasks).Error
	return tasks, err
}

func (r *TaskRepository) GetByStatus(status model.TaskStatus) ([]model.Task, error) {
	var tasks []model.Task
	err := DB.Where("status = ?", status).Order("created_at desc").Find(&tasks).Error
	return tasks, err
}

// GetByURL 未找到时返回 (nil, gorm.ErrRecordNotFound)，调用方用 errors.Is 区分。
func (r *TaskRepository) GetByURL(url string) (*model.Task, error) {
	var task model.Task
	if err := DB.Where("url = ?", url).First(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

// IsNotFound 用于调用方判断 GetByURL/GetByID 是否为"未找到"。
func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func (r *TaskRepository) UpdateNameAndHeader(id uint, name, header string) error {
	return DB.Model(&model.Task{}).
		Where("id = ?", id).
		Updates(map[string]any{"name": name, "header": header}).Error
}
