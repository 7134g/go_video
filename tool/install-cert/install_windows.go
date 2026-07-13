//go:build windows

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sys/windows"
)

func ensurePrivileged() {
	if isAdmin() {
		return
	}
	log.Println("需要管理员权限,正在请求提升...")
	runAsAdmin()
	os.Exit(0)
}

func isAdmin() bool {
	var sid *windows.SID
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
	verbPtr, _ := windows.UTF16PtrFromString("runas")
	exe, _ := os.Executable()
	exePtr, _ := windows.UTF16PtrFromString(exe)

	cwd, _ := os.Getwd()
	cwdPtr, _ := windows.UTF16PtrFromString(cwd)
	argPtr, _ := windows.UTF16PtrFromString(strings.Join(os.Args[1:], " "))

	var showCmd int32 = 1 // SW_SHOWNORMAL

	err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		fmt.Printf("提权失败: %v\n", err)
		os.Exit(1)
	}
}

func installCert(certPath string) error {
	cmd := exec.Command("certutil", "-addstore", "-f", "Root", certPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
