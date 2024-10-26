package io

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/TheTNB/panel/pkg/shell"
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
	cmd := exec.Command("chmod", "-R", fmt.Sprintf("%o", permission), path)
	return cmd.Run()
}

// Chown 修改文件或目录所有者
func Chown(path, user, group string) error {
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
	if err := os.Rename(src, dst); err != nil {
		// 在不同的文件系统中无法使用 os.Rename
		if _, err = shell.Execf(`mv -f '%s' '%s'`, src, dst); err != nil {
			return err
		}
	}

	return nil
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

// ReadDir 读取目录
func ReadDir(path string) ([]os.DirEntry, error) {
	return os.ReadDir(path)
}

// IsDir 判断是否为目录
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// SizeX 获取路径大小（du命令）
func SizeX(path string) (int64, error) {
	out, err := exec.Command("du", "-sb", path).Output()
	if err != nil {
		return 0, err
	}

	parts := strings.Fields(string(out))
	if len(parts) == 0 {
		return 0, fmt.Errorf("无法解析 du 输出")
	}

	return strconv.ParseInt(parts[0], 10, 64)
}

// CountX 统计目录下文件数
func CountX(path string) (int64, error) {
	out, err := exec.Command("find", path, "-printf", ".").Output()
	if err != nil {
		return 0, err
	}

	count := len(string(out))
	return int64(count), nil
}

// Search 查找文件/文件夹
func Search(path, keyword string, sub bool) (map[string]os.FileInfo, error) {
	paths := make(map[string]os.FileInfo)
	baseDepth := strings.Count(filepath.Clean(path), string(os.PathSeparator))

	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !sub && strings.Count(p, string(os.PathSeparator)) > baseDepth+1 {
			return filepath.SkipDir
		}
		if strings.Contains(info.Name(), keyword) {
			paths[p] = info
		}
		return nil
	})

	return paths, err
}

// SearchX 查找文件/文件夹（find命令）
func SearchX(path, keyword string, sub bool) (map[string]os.FileInfo, error) {
	paths := make(map[string]os.FileInfo)

	var cmd *exec.Cmd
	if sub {
		cmd = exec.Command("find", path, "-name", "*"+keyword+"*")
	} else {
		cmd = exec.Command("find", path, "-maxdepth", "1", "-name", "*"+keyword+"*")
	}
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		info, err := os.Stat(line)
		if err != nil {
			return nil, err
		}
		paths[line] = info
	}

	return paths, nil
}
