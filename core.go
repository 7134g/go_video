package main

import (
	"dv/base"
	"dv/config"
	"dv/table"
	"log"
	"sync"
)

type Core struct {
	base.Logger

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

func (c *Core) Run(tl []Cell) {
	for _, t := range tl {
		c.Submit(&t)
	}
}

func (c *Core) Wait() {
	c.wg.Wait()
	log.Println("本次运行结束")
}

func (c *Core) Submit(t *Cell) {
	c.wg.Add(1)
	c.vacancy <- struct{}{}
	go func() {
		defer c.wg.Done()
		err := t.Do()
		<-c.vacancy
		if err != nil {
			table.IncreaseErrorCount(t.Link)
			if table.IsMaxErrorCount(t.Link) {
				log.Printf("文件名：%v 超出最大尝试次数\n", t.Name)
				return
			}

			c.Submit(t)
		}

		// 完成
		c.doneCount++

	}()
}

func (c *Core) SetGroupCount(n int) {
	c.groupCount = n
}
