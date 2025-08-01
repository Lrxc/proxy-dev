package assets

import (
	"embed"
	"os"
	"path/filepath"
)

//go:embed *
var Resources embed.FS

func ReadByte(scriptPath string) ([]byte, error) {
	return Resources.ReadFile(scriptPath)
}

func Read(scriptPath string) []byte {
	file, err := Resources.ReadFile(scriptPath)
	if err != nil {
		return nil
	}
	return file
}

// 读取内嵌文件
func ReadFile(scriptPath string) (string, error) {
	//读取嵌入的脚本
	scriptByte, err := Resources.ReadFile(scriptPath)
	if err != nil {
		return "", err
	}

	filePath := filepath.Join("data", scriptPath)
	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}

	err = os.WriteFile(filePath, scriptByte, os.ModePerm)
	// 创建并写入文件
	if err != nil {
		return "", err
	}
	return filePath, nil
}
