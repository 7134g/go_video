package base

import (
	"log"
	"time"
)

var (
	interval       = 1              // 消息间隔时间（10s）
	errorMessage   = "[x] %5s %v\n" // 下载过程中错误消息
	workingMessage = "[√] %5s %v\n" // 下载过程中工作消息
	doneMessage    = "[O] %5s %v\n" // 下载完成中工作消息
)

type Logger struct {
}

func (l Logger) Fail(taskName, msg string) {
	log.Printf(errorMessage, taskName, msg)
}

func (l Logger) Doing(taskName, msg string) {
	log.Printf(workingMessage, taskName, msg)
}

func (l Logger) Done(taskName, msg string) {
	log.Printf(doneMessage, taskName, msg)
}

func (l Logger) Ticker() *time.Ticker {
	return time.NewTicker(time.Second * time.Duration(interval))
}

func SetInterval(value int) {
	interval = value
}

func GetInterval() int {
	return interval
}
