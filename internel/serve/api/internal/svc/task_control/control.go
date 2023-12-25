package task_control

import (
	"context"
	"dv/internel/serve/api/internal/util/calc"
	"dv/internel/serve/api/internal/util/model"
	"dv/internel/serve/api/internal/util/table"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"log"
	"strconv"
	"strings"
	"time"
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
	close(c.printStop)
	close(c.vacancy)
	c.vacancy = make(chan struct{}, len(c.vacancy))
	c.running = false
}

func (c *TaskControl) start() {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.running = true
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.printStop = make(chan struct{})
}

func (c *TaskControl) incDoneCount() {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.doneCount++
}

func (c *TaskControl) printDownloadProgress(taskTotal uint) {
	ticker := time.NewTicker(time.Second * 3)
	var lastDownloadTimeSince uint
	for {
		select {
		case <-c.printStop:
			return
		case <-ticker.C:
			nowDownloadDataLen, exist := table.M3u8DownloadDataLen.Get(c.Name)
			if !exist || nowDownloadDataLen == 0 {
				continue
			}

			downloadTimeSince := nowDownloadDataLen - lastDownloadTimeSince
			speed, unit := calc.UnitReturn(float64(downloadTimeSince))
			log.Println(fmt.Sprintf("%s 下载进度(%d/%d) 速度：%.2f %s/s 完成度：%.2f ",
				c.Name,
				c.doneCount, taskTotal,
				speed, unit,
				float64(c.doneCount)/float64(taskTotal)*100,
			) + "%")
			lastDownloadTimeSince = nowDownloadDataLen
		}
	}
}

func (c *TaskControl) submit(fn func() error, d *download) {
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

		keyPart := strings.Split(d.key, "_")
		taskId, _ := strconv.Atoi(keyPart[0])

		if err := fn(); err != nil {
			table.IncErrCount(d.key)
			if table.GetErrCount(d.key) >= tcConfig.TaskErrorMaxCount {
				_ = tasKModel.UpdateStatus(uint(taskId), model.StatusError)
				logx.Error(saveErrorCellData(d))
				return
			} else {
				logx.Error(d.key, err)
				c.submit(fn, d)
			}
		}

		c.incDoneCount()
		if len(keyPart) <= 2 {
			_ = tasKModel.UpdateStatus(uint(taskId), model.StatusSuccess)
			logx.Info(d.key, " is done")
		}

		return
	})
}

func (c *TaskControl) Run(task []model.Task) {
	defer c.Stop()
	c.start()

	go c.printDownloadProgress(uint(len(task)))

	for _, m := range task {
		w := newWork(m)
		d, particle := w.parseTask()
		if particle == nil {
			continue
		}

		_ = tasKModel.UpdateStatus(m.ID, model.StatusRunning)
		c.submit(particle, d)
	}
	c.wg.Wait()

}
