package main

import (
	"dv/config"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	defer func() {
		if err := any(recover()); err != nil {
			log.Println(err)
		}
		log.Println("按任意键结束……")
		var input string
		_, _ = fmt.Scanf("%v", &input)
	}()

	config.LoadConfig()

	core := NewCore()
	tl, err := ParseTaskList()
	if err != nil {
		log.Println("解析任务清单失败：", err)
		return
	}

	core.Run(tl)
	core.Wait()
	StoreHistory()
}

func StoreHistory() {
	if !config.GetConfig().TaskClear {
		return
	}
	newName := time.Now().Format("2006-01-02-15-04-05")
	_ = os.Rename(config.GetConfig().TaskList,
		fmt.Sprintf("%s/历史任务_%s.txt", config.GetConfig().HistoryDir, newName))
	// 清空任务清单
	f, err := os.Create(config.GetConfig().TaskList)
	if err != nil {
		log.Fatalln(err)
	}
	_ = f.Close()
}
