package repository

import (
	"luka/internal/model"

	"github.com/glebarez/sqlite"
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
