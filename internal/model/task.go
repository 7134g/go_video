package model

import "time"

type TaskStatus int

const (
	TaskStatusPending   TaskStatus = 0 // 待执行
	TaskStatusRunning   TaskStatus = 1 // 执行中
	TaskStatusCompleted TaskStatus = 2 // 完成
	TaskStatusFailed    TaskStatus = 3 // 失败
	TaskStatusPaused    TaskStatus = 4 // 已暂停
)

type Task struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Name      string     `json:"name" gorm:"not null"`
	URL       string     `json:"url" gorm:"not null"`
	Header    string     `json:"header"`
	Type      string     `json:"type" gorm:"not null"` // mp4 or m3u8
	Status    TaskStatus `json:"status" gorm:"default:0"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
