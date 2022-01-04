package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetRootDir 获取执行路径
func GetRootDir() string {
	// 文件不存在获取执行路径
	file, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		file = fmt.Sprintf(".%s", string(os.PathSeparator))
	} else {
		file = fmt.Sprintf("%s%s", file, string(os.PathSeparator))
	}
	return file
}

// GetCurrentPath 获取当前项目路径
func GetCurrentPath() string {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("获取当前项目路径err=", err)
	}
	return strings.Replace(path, "\\", "/", -1)
	return path
}

// PathExists 判断文件或目录是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
