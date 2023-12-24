package task_control

import (
	"context"
	"dv/internel/serve/api/internal/util/model"
	"dv/internel/serve/api/internal/util/table"
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
	close(c.vacancy)
	c.vacancy = make(chan struct{}, len(c.vacancy))
	c.running = false
}

func (c *TaskControl) start() {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.running = true
	c.ctx, c.cancel = context.WithCancel(context.Background())
}

func (c *TaskControl) incDoneCount() {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.doneCount++
}

func (c *TaskControl) Run(task []model.Task) {
	defer c.Stop()
	c.start()

	for _, m := range task {
		w := newWork(m)
		d, particle := w.parseTask()
		if particle == nil {
			continue
		}

		c.Submit(particle, d)
	}
	c.wg.Wait()
}

func (c *TaskControl) Submit(fn func() error, d *download) {
	c.wg.Add(1)
	select {
	case c.vacancy <- struct{}{}:
	case <-c.ctx.Done():
		logx.Info("cancel stop")
		return
	}
	go threading.GoSafe(func() {
		defer func() {
			c.wg.Done()
			<-c.vacancy
		}()

		if err := fn(); err != nil {
			table.IncErrCount(d.key)
			if table.GetErrCount(d.key) >= tcConfig.TaskErrorMaxCount {
				logx.Error(saveErrorCellData(d))
				return
			} else {
				logx.Error(d.key, err)
				c.Submit(fn, d)
			}
		}

		c.incDoneCount()
		logx.Info(d.key, "is done")

		return
	})
}
