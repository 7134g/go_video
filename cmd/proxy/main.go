package main

import (
	"fmt"
	"go_video/pkg/proxy"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"golang.org/x/sys/windows"
)

func main() {
	// 检查并要求管理员权限
	if !isAdmin() {
		log.Println("需要管理员权限，正在请求提升...")
		runAsAdmin()
		return
	}

	// 生成CA证书（首次运行）
	if err := proxy.GenCA(); err != nil {
		log.Printf("生成证书失败（可能已存在）: %v", err)
	}

	_ = InstallCert()
}

func isAdmin() bool {
	var sid *windows.SID
	// 虽然可以用打开物理驱动器的方法，但检查 WellKnownSid 更标准
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid,
	)
	if err != nil {
		return false
	}
	defer windows.FreeSid(sid)

	token := windows.Token(0)
	member, err := token.IsMember(sid)
	if err != nil {
		return false
	}
	return member
}

func runAsAdmin() {
	if isAdmin() {
		return
	}

	verbPtr, _ := windows.UTF16PtrFromString("runas")
	exe, _ := os.Executable()
	exePtr, _ := windows.UTF16PtrFromString(exe)

	// 拼接命令行参数，确保提权后的进程行为一致
	cwd, _ := os.Getwd()
	cwdPtr, _ := windows.UTF16PtrFromString(cwd)
	argPtr, _ := windows.UTF16PtrFromString(strings.Join(os.Args[1:], " "))

	var showCmd int32 = 1 // SW_SHOWNORMAL

	err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		fmt.Printf("提权失败: %v\n", err)
		os.Exit(1)
	}

	// 关键：退出当前无权限的进程
	os.Exit(0)
}

func InstallCert() error {
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
