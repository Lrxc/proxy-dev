//go:build windows
// +build windows

package util

import (
	"os/exec"
	"syscall"
)

func Exec(name string, arg ...string) ([]byte, error) {
	c := exec.Command(name, arg...)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return c.CombinedOutput()
}
