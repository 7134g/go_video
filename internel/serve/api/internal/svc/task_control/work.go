package task_control

import "dv/internel/serve/api/internal/model"

type work struct {
	task model.Task
}

func newWork(task model.Task) *work {
	return &work{task: task}
}
