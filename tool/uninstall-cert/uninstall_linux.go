//go:build linux

package main

import (
	"fmt"
	"os"
	"os/exec"
)

const linuxCertName = "proxy-ca.crt"

func ensurePrivileged() {}

func uninstallCert(cn string) error {
	removed := false
	for _, path := range []string{
		"/usr/local/share/ca-certificates/" + linuxCertName,
		"/etc/pki/ca-trust/source/anchors/" + linuxCertName,
	} {
		if err := os.Remove(path); err == nil {
			removed = true
		}
	}
	if !removed {
		return fmt.Errorf("未找到已安装的证书文件")
	}
	if _, err := exec.LookPath("update-ca-certificates"); err == nil {
		cmd := exec.Command("update-ca-certificates")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	if _, err := exec.LookPath("update-ca-trust"); err == nil {
		cmd := exec.Command("update-ca-trust", "extract")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	return nil
}
