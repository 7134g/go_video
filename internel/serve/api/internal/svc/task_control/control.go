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
	"path/filepath"
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

func (c *TaskControl) printDownloadProgress(name string, taskTotal uint) {
	if name == "" {
		return
	}

	ticker := time.NewTicker(time.Second * 3)
	var lastDownloadTimeSince uint
	for {
		select {
		case <-c.printStop:
			return
		case <-ticker.C:
			nowDownloadDataLen, exist := table.DownloadDataLen.Get(name)
			if !exist {
				continue
			}

			downloadTimeSince := nowDownloadDataLen - lastDownloadTimeSince
			speed, unit := calc.UnitReturn(float64(downloadTimeSince))
			log.Println(fmt.Sprintf("%s 下载进度(%d/%d) 速度：%.2f %s/s 完成度：%.2f ",
				name,
				c.doneCount, taskTotal,
				speed/3, unit,
				float64(c.doneCount)/float64(taskTotal)*100,
			) + "%")
			lastDownloadTimeSince = nowDownloadDataLen
		}
	}
}

func (c *TaskControl) submit(fn particleFunc, params []any) {
	c.wg.Add(1)
	select {
	case c.vacancy <- struct{}{}:
	case <-c.ctx.Done():
		logx.Info("cancel stop")
		return
	}

	d := params[0].(*download)
	go threading.GoSafe(func() {
		if d == nil {
			return
		}
		defer func() {
			c.wg.Done()
			<-c.vacancy
		}()

		keyPart := strings.Split(d.key, "_")
		taskId, _ := strconv.Atoi(keyPart[0])

		if err := fn([]any{d}); err != nil {
			table.IncErrCount(d.key)
			if table.GetErrCount(d.key) >= tcConfig.TaskErrorMaxCount {
				_ = tasKModel.UpdateStatus(uint(taskId), model.StatusError)
				logx.Error(keyPart[1], "任务失败")
			} else {
				logx.Errorw(
					"error message",
					logx.Field("retry_count", table.GetErrCount(d.key)),
					logx.Field("key", d.key),
					logx.Field("error", err),
				)
				time.Sleep(time.Second * time.Duration(tcConfig.TaskErrorDuration))
				c.submit(fn, []any{d}) // 重试
			}
			return
		}

		c.incDoneCount()
		if len(keyPart) <= 2 {
			_ = tasKModel.UpdateStatus(uint(taskId), model.StatusSuccess)
			logx.Info(d.key, " is done")
		}

		return
	})
}

func (c *TaskControl) Run(tasks []model.Task) {
	logx.Info("running ......")
	logx.Info(filepath.Abs(tcConfig.SaveDir))
	defer c.Stop()
	c.start()

	//go c.printDownloadProgress(uint(len(tasks)))

	for _, m := range tasks {
		w := newWork(m)
		particle, d := w.parseTask()
		if particle == nil {
			continue
		}

		_ = tasKModel.UpdateStatus(m.ID, model.StatusRunning)
		logx.Infof("=========> 任务开始：%s", w.task.Name)
		c.submit(particle, []any{d})
	}
	c.wg.Wait()

	logx.Info("所有任务已经结束 <=========")
}
