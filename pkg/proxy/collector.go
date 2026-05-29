package proxy

import (
	"fmt"
	"net/http"
)

// Collector 把 MITM 嗅探到的视频任务投递给 service 层。
// tasks channel 缓冲 100；下游 doTask 处理慢于上游捕获时 Collect 会阻塞，
// 进而反压住 ModifyResponse —— 故意如此，避免内存里堆积无界任务。
type Collector struct {
	tasks chan *VideoTask
}

func NewCollector() *Collector {
	return &Collector{
		tasks: make(chan *VideoTask, 100),
	}
}

func (c *Collector) Collect(req *http.Request, title, videoType string) {
	fmt.Println("抓取到新任务: ", title, req.URL.String())
	task := Capture(req)
	task.Type = videoType
	if title != "" {
		task.Title = title
	}
	c.tasks <- task
}

func (c *Collector) Tasks() <-chan *VideoTask {
	return c.tasks
}

func (c *Collector) Close() {
	close(c.tasks)
}
