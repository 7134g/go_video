package task_control

import (
	"dv/internel/serve/api/internal/model"
	"encoding/json"
	"time"
)

func saveErrorCellData(c *cell) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return errModel.Insert(&model.Error{
		TaskId:     c.TaskId,
		Data:       string(b),
		CreateTime: time.Now(),
	})
}
