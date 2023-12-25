package model

import (
	"errors"
	"gorm.io/gorm"
)

type Task struct {
	ID        uint   `json:"id" gorm:"primaryKey;column:id"`
	Name      string `json:"name"`       // 任务名字
	VideoType string `json:"video_type"` // 视频类型
	Type      string `json:"type"`       // 任务类型
	Data      string `json:"data"`       // url 或者 curl
	Status    uint   `json:"status"`     // 执行状态
}

const (
	VideoTypeMp4  = "mp4"
	VideoTypeM3u8 = "m3u8"

	TypeUrl  = "url"
	TypeCurl = "curl"
)

type statusEnum uint

const (
	StatsuWait statusEnum = iota
	StatsuRunning
	StatsuError
	StatsuSuccess
)

type TaskModel struct {
	DB *gorm.DB
}

func NewTaskModel(db *gorm.DB) *TaskModel {
	return &TaskModel{DB: db}
}

func (m *TaskModel) List() ([]Task, error) {
	var data []Task
	err := m.DB.Model(&Task{}).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (m *TaskModel) Update(task *Task) error {
	return m.DB.Model(&Task{}).Where("id = ?", task.ID).Updates(task).Error
}

func (m *TaskModel) UpdateStatus(id, status uint) error {
	switch statusEnum(status) {
	case StatsuWait:
	case StatsuRunning:
	case StatsuError:
	case StatsuSuccess:
	default:
		return errors.New("status error")
	}

	return m.DB.Model(&Task{}).Where("id = ?", id).Update("status= ?", status).Error
}

func (m *TaskModel) Exist(data string) (*Task, error) {
	findTask := &Task{}
	if err := m.DB.Model(&Task{}).Where("data = ?", data).First(findTask).Error; err != nil {
		return nil, err
	}

	return findTask, nil
}

func (m *TaskModel) Insert(task *Task) error {
	findTask := Task{}
	err := m.DB.Model(&Task{}).Where("data = ?", task.Data).First(&findTask).Error
	switch err {
	case gorm.ErrRecordNotFound:
		return m.DB.Model(&Task{}).Create(task).Error
	case nil:
		*task = findTask
		return nil
	default:
		return err
	}

	//return m.DB.Model(&Task{}).Create(task).Error
}

func (m *TaskModel) Delete(id uint) error {
	return m.DB.Model(&Task{}).Delete("id = ?", id).Error
}
