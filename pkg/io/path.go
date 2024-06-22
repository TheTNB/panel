package io

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/goravel/framework/support/env"
)

// Remove 删除文件/目录
func Remove(path string) error {
	return os.RemoveAll(path)
}

// Mkdir 创建目录
func Mkdir(path string, permission os.FileMode) error {
	return os.MkdirAll(path, permission)
}

// Chmod 修改文件/目录权限
func Chmod(path string, permission os.FileMode) error {
	if env.IsWindows() {
		return errors.New("chmod is not supported on Windows")
	}

	cmd := exec.Command("chmod", "-R", fmt.Sprintf("%o", permission), path)
	return cmd.Run()
}

// Chown 修改文件或目录所有者
func Chown(path, user, group string) error {
	if env.IsWindows() {
		return errors.New("chown is not supported on Windows")
	}

	cmd := exec.Command("chown", "-R", user+":"+group, path)
	return cmd.Run()
}

// Exists 判断路径是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Empty 判断路径是否为空
func Empty(path string) bool {
	files, err := os.ReadDir(path)
	if err != nil {
		return true
	}

	return len(files) == 0
}

func Mv(src, dst string) error {
	err := os.Rename(src, dst)
	if err != nil {
		// 如果在不同的文件系统中移动文件，os.Rename 可能会失败
		err = Cp(src, dst)
		if err != nil {
			return err
		}
		err = os.RemoveAll(src)
	}

	return err
}

// Cp 复制文件或目录
func Cp(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if srcInfo.IsDir() {
		return copyDir(src, dst)
	}
	return copyFile(src, dst)
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = copyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			err = copyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Size 获取路径大小
func Size(path string) (int64, error) {
	var size int64

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		size += info.Size()
		return nil
	})

	return size, err
}

// TempDir 创建临时目录
func TempDir(prefix string) (string, error) {
	return os.MkdirTemp("", prefix)
}
