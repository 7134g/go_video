package task_control

import (
	"dv/internel/serve/api/internal/model"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

func saveErrorCellData(d *download) error {
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}

	taskId := strings.Split(d.key, "_")
	if len(taskId) != 2 {
		return errors.New("key error")
	}
	id, err := strconv.Atoi(taskId[0])
	if err != nil {
		return err
	}
	return errModel.Insert(&model.Error{
		TaskId:     uint(id),
		Data:       string(b),
		CreateTime: time.Now(),
	})
}
