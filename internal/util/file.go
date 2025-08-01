package util

import (
	"os"
	"path/filepath"
)

// 创建父级路径
func CreateParentFile(filename string) {
	//当前绝对路径
	path, _ := filepath.Abs(filename)
	//上级路径
	path = filepath.Join(path, "../")
	//创建目录
	err := os.MkdirAll(path, 0666)

	if err != nil {
		panic(err)
	}
}

func FileExist(filename string) bool {
	abs, _ := filepath.Abs(filename)

	_, err := os.Stat(abs)
	if err == nil {
		return true // 文件存在
	}
	if os.IsNotExist(err) {
		return false // 文件不存在
	}
	return false // 其他错误（如权限不足），通常也视为不存在
}
