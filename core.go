package main

import (
	"dv/base"
	"dv/config"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Core struct {
	base.Logger

	wg        *sync.WaitGroup
	list      []*Task       // 所有任务
	retryTask chan *Task    // 需要重试的任务
	vacancy   chan struct{} // 空缺

	doneCount  int // 已完成任务数
	groupCount int // 该组任务总数量

	cfg *config.ProjectConfig
}

func NewCore(cfg *config.ProjectConfig) Core {
	return Core{
		wg:        &sync.WaitGroup{},
		list:      make([]*Task, 0),
		retryTask: make(chan *Task, 1000),
		vacancy:   make(chan struct{}, cfg.Concurrency),
		cfg:       cfg,
	}
}

func (c *Core) Run() {
	if c.list == nil || len(c.list) == 0 {
		return
	}
	for _, task := range c.list {
		c.Submit(task)
	}
	// 等待任务完成
	c.wg.Wait()

	// 下面是重试
	if len(c.retryTask) == 0 {
		// 没有失败任务
		return
	}
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
	for {
		select {
		case t := <-c.retryTask:
			c.Submit(t)
			if t.errorCount >= config.GetConfig().TaskErrorMaxCount {
				log.Printf("文件名：%v 超出最大尝试次数\n", t.fileName)
				continue
			}
		case <-ticker.C:
			if len(c.vacancy) == 0 {
				return
			}
		}
	}
}

func (c *Core) Wait() {
	c.wg.Wait()
	log.Println("本次运行结束")
}

func (c *Core) Submit(t *Task) {
	c.wg.Add(1)
	c.vacancy <- struct{}{}
	go func() {
		t.Do() // 执行
		c.wg.Done()
		if t.errorCount > 0 {
			time.Sleep(time.Second * time.Duration(config.GetConfig().TaskErrorDuration))
			c.retryTask <- t
		} else {
			c.doneCount++
			c.printM3u8(t)
		}

		<-c.vacancy
	}()
}

func (c *Core) AddTask(t *Task) {
	c.list = append(c.list, t)
}

func (c *Core) SetGroup(n int) {
	c.groupCount = n
}

func (c *Core) ParseTaskList() error {
	bs, err := os.ReadFile(c.cfg.TaskName)
	if err != nil {
		return err
	}

	content := string(bs)
	if content == "" {
		return errors.New("content is 0")
	}
	reHead, _ := regexp.Compile(`\s+`)
	content = reHead.ReplaceAllString(content, "\n")
	content = strings.TrimPrefix(content, "\n")
	content = strings.TrimSuffix(content, "\n")
	list := strings.Split(content, "\n")
	for i := 0; i < len(list); i++ {
		if i+1 == len(list) {
			break
		}
		key := list[i]
		value := list[i+1]
		if len(value) < 4 || "http" != value[:4] {
			log.Println("错误值：", value)
			continue
		}
		task := NewTask(key, config.GetConfig().SaveDir, value)
		c.list = append(c.list, task)
		i++
	}

	return nil
}

func (c *Core) StoreHistory() {
	if !c.cfg.TaskClear {
		return
	}
	newName := time.Now().Format("2006-01-02-15-04-05")
	_ = os.Rename(c.cfg.TaskName, fmt.Sprintf("%s/历史任务_%s.txt", c.cfg.HistoryDir, newName))
	// 清空任务清单
	f, err := os.Create(c.cfg.TaskName)
	if err != nil {
		log.Fatalln(err)
	}
	_ = f.Close()
}

func (c *Core) printM3u8(t *Task) {
	if c.groupCount == 0 {
		return
	}
	// m3u8组任务打印信息
	c.Doing(t.fileName, fmt.Sprintf("分片下载进度(%d/%d) %.2f ",
		c.doneCount, c.groupCount, float64(c.doneCount)*100/float64(c.groupCount))+"%")
}
