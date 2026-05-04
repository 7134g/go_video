package proxy

import (
	"fmt"
	"net/http"
)

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
