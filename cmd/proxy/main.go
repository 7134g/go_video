package main

import (
	"fmt"
	"log"

	"go_video/pkg/proxy"
)

func main() {
	// 生成CA证书（首次运行）
	if err := proxy.GenCA("ca.crt", "ca.key"); err != nil {
		log.Printf("生成证书失败（可能已存在）: %v", err)
	}

	// 创建代理服务器
	server, err := proxy.NewServer("ca.crt", "ca.key")
	if err != nil {
		log.Fatal(err)
	}

	// 启动任务收集
	go func() {
		for task := range server.Tasks() {
			fmt.Printf("捕获视频: %s\n", task.URL)
		}
	}()

	// 启动代理
	fmt.Println("代理服务器启动在 127.0.0.1:10888")
	if err := server.Listen("127.0.0.1:10888"); err != nil {
		log.Fatal(err)
	}
}
