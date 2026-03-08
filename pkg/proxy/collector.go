package proxy

type Collector struct {
	tasks chan *VideoTask
}

func NewCollector() *Collector {
	return &Collector{
		tasks: make(chan *VideoTask, 100),
	}
}

func (c *Collector) Collect(task *VideoTask) {
	if task != nil && task.Title != "" {
		c.tasks <- task
	}
}

func (c *Collector) Tasks() <-chan *VideoTask {
	return c.tasks
}

func (c *Collector) Close() {
	close(c.tasks)
}
