package io

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver/v3"
)

// Write 写入文件
func Write(path string, data string, permission os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), permission); err != nil {
		return err
	}

	err := os.WriteFile(path, []byte(data), permission)
	if err != nil {
		return err
	}

	return nil
}

// WriteAppend 追加写入文件
func WriteAppend(path string, data string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}

// Read 读取文件
func Read(path string) (string, error) {
	data, err := os.ReadFile(path)
	return string(data), err
}

// FileInfo 获取文件大小
func FileInfo(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

// UnArchive 智能解压文件
func UnArchive(file string, dst string) error {
	return archiver.Unarchive(file, dst)
}

// Archive 智能压缩文件
func Archive(src []string, dst string) error {
	return archiver.Archive(src, dst)
}

// TempFile 创建临时文件
func TempFile(prefix string) (*os.File, error) {
	return os.CreateTemp("", prefix)
}

// IsSymlink 判读是否为软链接
func IsSymlink(mode os.FileMode) bool {
	return mode&os.ModeSymlink != 0
}

// IsHidden 判断是否为隐藏文件
func IsHidden(path string) bool {
	_, file := filepath.Split(path)
	return strings.HasPrefix(file, ".")
}

// GetSymlink 获取软链接目标
func GetSymlink(path string) string {
	linkPath, err := os.Readlink(path)
	if err != nil {
		return ""
	}
	return linkPath
}
