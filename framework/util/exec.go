package util

import "os"

// 获取当前执行程序目录
func GetExecDirectory() string {
	file, err := os.Getwd()
	if err != nil {
		return ""
	}
	return file + "/"
}
