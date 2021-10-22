package util

import (
	"os"
	"syscall"
)

// 获取当前执行程序目录
func GetExecDirectory() string {
	file, err := os.Getwd()
	if err != nil {
		return ""
	}
	return file + "/"
}

func CheckProcessExist(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = process.Signal(syscall.Signal(0))
	if err != nil {
		return false
	}
	return true
}
