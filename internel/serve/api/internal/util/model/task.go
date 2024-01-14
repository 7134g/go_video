package model

import (
	"dv/internel/serve/api/internal/types"
	"errors"
	"fmt"
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
)

const (
	TypeUrl   = "url"
	TypeCurl  = "curl"
	TypeProxy = "proxy"
)

type statusEnum uint

const (
	StatusWait statusEnum = iota
	StatusRunning
	StatusError
	StatusSuccess
)

type TaskModel struct {
	DB *gorm.DB
}

func NewTaskModel(db *gorm.DB) *TaskModel {
	return &TaskModel{DB: db}
}

func (m *TaskModel) parseWhere(where map[string]any) *gorm.DB {
	_db := m.DB.Model(&Task{})
	for key, value := range where {
		if value == nil {
			continue
		}
		switch key {
		case "type":
			if value == "all" {
				continue
			}
		case "video_type":

		}
		sql := fmt.Sprintf("%s = ?", key)
		_db = _db.Where(sql, value)
	}

	return _db
}

func (m *TaskModel) Count(turner types.PageTurner) (int64, error) {
	where, _, _ := turner.ParseMysql()

	_db := m.parseWhere(where)
	var count int64
	if err := _db.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (m *TaskModel) List(turner types.PageTurner) ([]Task, error) {
	where, offset, limit := turner.ParseMysql()

	_db := m.parseWhere(where)
	var data []Task
	err := _db.Offset(offset).Limit(limit).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (m *TaskModel) Update(task *Task) error {
	return m.DB.Model(&Task{}).Where("id = ?", task.ID).Updates(task).Error
}

func (m *TaskModel) UpdateStatus(id uint, status statusEnum) error {
	switch status {
	case StatusWait:
	case StatusRunning:
	case StatusError:
	case StatusSuccess:
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
