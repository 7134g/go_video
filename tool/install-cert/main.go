package main

import (
	"go_video/pkg/proxy"
	"log"
	"os"
)

func main() {
	ensurePrivileged()

	if _, err := os.Stat(proxy.CACertFile); os.IsNotExist(err) {
		if err := proxy.GenCA(); err != nil {
			log.Fatal("生成 CA 证书失败: ", err)
		}
	}

	if err := installCert(proxy.CACertFile); err != nil {
		log.Fatal("安装证书失败: ", err)
	}
	log.Println("证书安装完成")
}
