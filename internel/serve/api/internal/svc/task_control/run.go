package task_control

import (
	"context"
	"dv/internel/serve/api/internal/model"
	"fmt"
)

func (t *TaskControl) GetStatus() bool {
	t.mux.Lock()
	defer t.mux.Unlock()

	return t.running
}

func (t *TaskControl) Stop() {
	t.mux.Lock()
	defer t.mux.Unlock()

	t.cancel()
	t.running = false
}

func (t *TaskControl) Run(task []model.Task) {
	t.ctx, t.cancel = context.WithCancel(context.Background())

	for _, m := range task {
		// todo
		w := newWork(m)
		fmt.Println(w)
	}
}
