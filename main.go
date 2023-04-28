package main

import (
	"dv/config"
	"fmt"
	"log"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
		log.Println("按任意键结束……")
		var input string
		_, _ = fmt.Scanf("%v", &input)
	}()

	config.LoadConfig()

	core := NewCore(config.GetConfig())
	if err := core.ParseTaskList(); err != nil {
		log.Println("解析任务清单失败：", err)
		return
	}

	core.Run()
	core.Wait()
	core.StoreHistory()
}
