package main

import (
	"go_video/pkg/proxy"
	"log"
	"os"
)

func main() {
	ensurePrivileged()

	if _, err := os.Stat(proxy.CACertFile); os.IsNotExist(err) {
		log.Println("CA 证书文件不存在,尝试从系统信任库直接卸载...")
	} else {
		ca, err := proxy.LoadCA()
		if err != nil {
			log.Fatal("读取 CA 证书失败: ", err)
		}
		if err := uninstallCert(ca.Subject.CommonName); err != nil {
			log.Fatal("卸载证书失败: ", err)
		}
	}

	removeIfExists(proxy.CACertFile)
	removeIfExists(proxy.CAKeyFile)
	log.Println("证书卸载完成")
}

func removeIfExists(path string) {
	if err := os.Remove(path); err == nil {
		log.Println("已删除:", path)
	}
}
