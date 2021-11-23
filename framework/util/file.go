package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断是否是隐藏目录
func IsHiddenDirectory(path string) bool {
	return len(path) > 1 && strings.HasPrefix(filepath.Base(path), ".")
}

// 返回当前目录的子目录的目录名
func SubDir(folder string) ([]string, error) {
	dirs, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, err
	}
	var subDirs []string
	for _, v := range dirs {
		if v.IsDir() {
			subDirs = append(subDirs, v.Name())
		}
	}
	return subDirs, nil
}
