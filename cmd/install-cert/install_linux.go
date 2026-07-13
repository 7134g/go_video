//go:build linux

package main

import (
	"flag"
	"fmt"
	"go_video/pkg/proxy"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

const linuxCertName = "proxy-ca.crt"

var systemOnly = flag.Bool("system-only", false, "internal: skip NSS, install system trust only (used after sudo re-exec)")

func ensurePrivileged() {
	flag.Parse()
}

func installCert(certPath string) error {
	if !*systemOnly && os.Geteuid() != 0 {
		if err := installNSS(certPath); err != nil {
			log.Printf("NSS 安装失败(不致命): %v", err)
		}
		return sudoReexec()
	}
	return installSystem(certPath)
}

func installNSS(certPath string) error {
	if _, err := exec.LookPath("certutil"); err != nil {
		fmt.Println("提示: 未检测到 certutil(libnss3-tools),Chrome/Chromium 信任 NSS 库,无法自动写入。")
		fmt.Println("      可执行: sudo apt install libnss3-tools  或  sudo dnf install nss-tools  后重跑")
		return nil
	}
	ca, err := proxy.LoadCA()
	if err != nil {
		return err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	nssDir := filepath.Join(home, ".pki", "nssdb")
	if err := os.MkdirAll(nssDir, 0755); err != nil {
		return err
	}
	cmd := exec.Command("certutil", "-d", "sql:"+nssDir, "-A",
		"-t", "C,,", "-n", ca.Subject.CommonName, "-i", certPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	fmt.Println("已写入 NSS 库:", nssDir)
	return nil
}

func sudoReexec() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	fmt.Println("接下来需要 sudo 安装系统根证书:")
	args := []string{"sudo", exe, "--system-only"}
	return syscall.Exec("/usr/bin/sudo", args, os.Environ())
}

func installSystem(certPath string) error {
	if _, err := exec.LookPath("update-ca-certificates"); err == nil {
		dst := "/usr/local/share/ca-certificates/" + linuxCertName
		if err := copyFile(certPath, dst, 0644); err != nil {
			return err
		}
		cmd := exec.Command("update-ca-certificates")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	if _, err := exec.LookPath("update-ca-trust"); err == nil {
		dst := "/etc/pki/ca-trust/source/anchors/" + linuxCertName
		if err := copyFile(certPath, dst, 0644); err != nil {
			return err
		}
		cmd := exec.Command("update-ca-trust", "extract")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	return fmt.Errorf("未找到 update-ca-certificates 或 update-ca-trust,请手动安装证书到系统信任库")
}

func copyFile(src, dst string, perm os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
