package utils

import (
	"os"
)

// IsWritableDirectory 检查目录是否能够写入
func IsWritableDirectory(path string) bool {
	// 尝试在目录中创建临时文件
	tmpFile, err := os.CreateTemp(path, "test-write")
	if err != nil {
		// 创建文件失败，目录可能不可写
		return false
	}

	// 清理：删除创建的临时文件
	defer os.Remove(tmpFile.Name())
	return true
}
