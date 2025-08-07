//go:build !windows
// +build !windows

package util

import (
	"os/exec"
)

func Exec(name string, arg ...string) ([]byte, error) {
	c := exec.Command(name, arg...)
	return c.CombinedOutput()
}
