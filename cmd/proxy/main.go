package main

import (
	"log"
	"os/exec"
	"runtime"

	"go_video/pkg/proxy"
)

func main() {
	// 生成CA证书（首次运行）
	if err := proxy.GenCA(); err != nil {
		log.Printf("生成证书失败（可能已存在）: %v", err)
	}

	_ = InstallCert()
}

func InstallCert() error {
	// certutil -addstore -f "Root" "你的证书文件路径.cer"
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("certutil", "-addstore", "-f", "Root", proxy.CACertFile)
		return cmd.Run()
	case "darwin":
		cmd := exec.Command("sudo", "security", "add-trusted-cert", "-d", "-r", "trustRoot", "-k", "/Library/Keychains/System.keychain", proxy.CACertFile)
		return cmd.Run()
	default:
		return nil
	}
}
