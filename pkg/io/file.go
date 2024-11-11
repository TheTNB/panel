package io

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/TheTNB/panel/pkg/chattr"
)

// Write 写入文件
func Write(path string, data string, permission os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), permission); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_RDONLY, permission)
	if err != nil {
		return err
	}

	iFlag, _ := chattr.IsAttr(file, chattr.FS_IMMUTABLE_FL)
	aFlag, _ := chattr.IsAttr(file, chattr.FS_APPEND_FL)
	if iFlag {
		_ = chattr.UnsetAttr(file, chattr.FS_IMMUTABLE_FL)
	}
	if aFlag {
		_ = chattr.UnsetAttr(file, chattr.FS_APPEND_FL)
	}

	// 关闭文件重新以写入方式打开
	if err = file.Close(); err != nil {
		return err
	}
	file, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, permission)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		return err
	}

	if iFlag {
		_ = chattr.SetAttr(file, chattr.FS_IMMUTABLE_FL)
	}
	if aFlag {
		_ = chattr.SetAttr(file, chattr.FS_APPEND_FL)
	}

	return nil
}

// WriteAppend 追加写入文件
func WriteAppend(path string, data string, permission os.FileMode) error {
	file, err := os.OpenFile(path, os.O_RDONLY, permission)
	if err != nil {
		return err
	}

	iFlag, _ := chattr.IsAttr(file, chattr.FS_IMMUTABLE_FL)
	if iFlag {
		_ = chattr.UnsetAttr(file, chattr.FS_IMMUTABLE_FL)
	}

	// 关闭文件重新以写入方式打开
	if err = file.Close(); err != nil {
		return err
	}
	file, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, permission)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		return err
	}

	if iFlag {
		_ = chattr.SetAttr(file, chattr.FS_IMMUTABLE_FL)
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

// TempFile 创建临时文件
func TempFile(dir, prefix string) (*os.File, error) {
	return os.CreateTemp(dir, prefix)
}
