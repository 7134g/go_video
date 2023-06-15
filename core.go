package main

import (
	"dv/config"
	"dv/table"
	"log"
	"sync"
	"time"
)

type Core struct {
	wg      *sync.WaitGroup
	vacancy chan struct{} // 并发控制

	doneCount  int // 已完成任务数
	groupCount int // 该组任务总数量

}

func NewCore() Core {
	return Core{
		wg:      &sync.WaitGroup{},
		vacancy: make(chan struct{}, config.GetConfig().Concurrency),
	}
}

func (c *Core) Run(ts []Task) {
	if ts == nil || len(ts) == 0 {
		return
	}

	for i := 0; i < len(ts); i++ {
		c.Submit(&ts[i])
	}

}

func (c *Core) Wait() {
	c.wg.Wait()
	//log.Println("本次运行结束")
	if len(table.RangeErrorTask()) != 0 {
		log.Println("失败任务：", table.RangeErrorTask())
	}
}

func (c *Core) Submit(t *Task) {
	c.wg.Add(1)
	c.vacancy <- struct{}{}
	go func() {
		defer c.wg.Done()
		err := t.Do()
		<-c.vacancy
		if err != nil {
			log.Println(t.Name, err)
			table.IncreaseErrorCount(t.Link)
			if table.IsMaxErrorCount(t.Link) {
				table.AddErrorTask(t.Name, t)
				log.Printf("文件名：%v 超出最大尝试次数\n", t.Name)
				return
			}
			time.Sleep(time.Second * time.Duration(config.GetConfig().TaskErrorDuration))
			c.Submit(t)
			return
		}

		// 完成
		c.doneCount++

	}()
}

func (c *Core) SetGroupCount(n int) {
	c.groupCount = n
}
