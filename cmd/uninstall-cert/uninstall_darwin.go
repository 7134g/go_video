//go:build darwin

package main

import (
	"os"
	"os/exec"
)

func ensurePrivileged() {}

func uninstallCert(cn string) error {
	cmd := exec.Command("sudo", "security", "delete-certificate", "-c", cn,
		"/Library/Keychains/System.keychain")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
