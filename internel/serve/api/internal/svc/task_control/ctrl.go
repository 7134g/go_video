package task_control

import (
	"context"
	"dv/internel/serve/api/internal/model"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
)

func (c *TaskControl) GetStatus() bool {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.running
}

func (c *TaskControl) Stop() {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.cancel()
	c.running = false
}

func (c *TaskControl) Run(task []model.Task) {
	c.running = true
	c.ctx, c.cancel = context.WithCancel(context.Background())

	for _, m := range task {
		w := newWork(m)
		particle := w.parseTask()
		if particle == nil {
			continue
		}
		c.Submit(func() {
			if err := particle.Do(); err != nil {
				logx.Error(err)
			}
		})
	}
	c.wg.Wait()
	c.Stop()
}

func (c *TaskControl) Submit(fn func(), deriveFlag ...bool) {
	c.wg.Add(1)
	if len(deriveFlag) == 0 {
		c.vacancy <- struct{}{}
	}
	go threading.GoSafe(func() {
		defer func() {
			c.wg.Done()
			if len(deriveFlag) == 0 {
				<-c.vacancy
			}
		}()
		fn()
	})
}
