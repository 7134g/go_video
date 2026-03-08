package repository

import (
	"go_video/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("video.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	return DB.AutoMigrate(&model.Task{})
}
