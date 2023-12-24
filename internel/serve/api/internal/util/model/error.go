package model

import (
	"gorm.io/gorm"
	"time"
)

type Error struct {
	ID         uint      `json:"id" gorm:"primaryKey;column:id"`
	TaskId     uint      `gorm:"column:task_id" json:"task_id"`
	Data       string    `gorm:"column:data" json:"data"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
}

type ErrorModel struct {
	DB *gorm.DB
}

func NewErrorModel(db *gorm.DB) *ErrorModel {
	return &ErrorModel{DB: db}
}

func (e *ErrorModel) List() ([]*Error, error) {
	data := make([]*Error, 0)
	err := e.DB.Model(&Error{}).Find(&data).Error

	return data, err
}

func (e *ErrorModel) Insert(data *Error) error {
	return e.DB.Model(&Error{}).Create(data).Error
}
