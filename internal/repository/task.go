package repository

import (
	"go_video/internal/model"
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
	err := DB.First(&task, id).Error
	return &task, err
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

func (r *TaskRepository) UpdateStatus(id uint, status model.TaskStatus) error {
	return DB.Model(&model.Task{}).Where("id = ?", id).Update("status", status).Error
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
