package tools

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/goravel/framework/facades"
)

// WriteFile 写入文件
func WriteFile(path string, data string, permission os.FileMode) bool {
	if err := os.MkdirAll(filepath.Dir(path), permission); err != nil {
		facades.Log().Errorf("[面板][Helpers] 创建目录失败: %s", err.Error())
		return false
	}

	err := os.WriteFile(path, []byte(data), permission)
	if err != nil {
		facades.Log().Errorf("[面板][Helpers] 写入文件 %s 失败: %s", path, err.Error())
		return false
	}

	return true
}

// ReadFile 读取文件
func ReadFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		facades.Log().Errorf("[面板][Helpers] 读取文件 %s 失败: %s", path, err.Error())
		return ""
	}

	return string(data)
}

// RemoveFile 删除文件
func RemoveFile(path string) bool {
	if err := os.Remove(path); err != nil {
		facades.Log().Errorf("[面板][Helpers] 删除文件 %s 失败: %s", path, err.Error())
		return false
	}

	return true
}

// ExecShell 执行 shell 命令
func ExecShell(shell string) string {
	cmd := exec.Command("bash", "-c", shell)

	output, err := cmd.CombinedOutput()
	if err != nil {
		facades.Log().Errorf("[面板][Helpers] 执行命令 %s 失败: %s", shell, err.Error())
		return ""
	}

	return strings.TrimSpace(string(output))
}

// ExecShellAsync 异步执行 shell 命令
func ExecShellAsync(shell string) {
	cmd := exec.Command("bash", "-c", shell)

	err := cmd.Start()
	if err != nil {
		facades.Log().Errorf("[面板][Helpers] 执行命令 %s 失败: %s", shell, err.Error())
	}

	go func() {
		err := cmd.Wait()
		if err != nil {
			facades.Log().Errorf("[面板][Helpers] 执行命令 %s 失败: %s", shell, err.Error())
		}
	}()
}

// Mkdir 创建目录
func Mkdir(path string, permission os.FileMode) bool {
	if err := os.MkdirAll(path, permission); err != nil {
		facades.Log().Errorf("[面板][Helpers] 创建目录 %s 失败: %s", path, err.Error())
		return false
	}

	return true
}

// Chmod 修改文件权限
func Chmod(path string, permission os.FileMode) bool {
	if err := os.Chmod(path, permission); err != nil {
		facades.Log().Errorf("[面板][Helpers] 修改文件 %s 权限失败: %s", path, err.Error())
		return false
	}

	return true
}

// Chown 修改路径所有者
func Chown(path, user, group string) bool {
	cmd := exec.Command("chown", "-R", user+":"+group, path)

	err := cmd.Run()
	if err != nil {
		facades.Log().Errorf("[面板][Helpers] 修改路径 %s 所有者失败: %s", path, err.Error())
		return false
	}

	return true
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

// Mv 移动路径
func Mv(src, dst string) bool {
	cmd := exec.Command("mv", src, dst)

	err := cmd.Run()
	if err != nil {
		facades.Log().Errorf("[面板][Helpers] 移动 %s 到 %s 失败: %s", src, dst, err.Error())
		return false
	}

	return true
}

// Cp 复制路径
func Cp(src, dst string) bool {
	cmd := exec.Command("cp", "-r", src, dst)

	err := cmd.Run()
	if err != nil {
		facades.Log().Errorf("[面板][Helpers] 复制 %s 到 %s 失败: %s", src, dst, err.Error())
		return false
	}

	return true
}

// Size 获取路径大小
func Size(path string) int64 {
	var size int64

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		size += info.Size()
		return nil
	})
	if err != nil {
		facades.Log().Errorf("[面板][Helpers] 获取路径 %s 大小失败: %s", path, err.Error())
		return 0
	}

	return size
}

// FileSize 获取文件大小
func FileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		facades.Log().Errorf("[面板][Helpers] 获取文件 %s 大小失败: %s", path, err.Error())
		return 0
	}

	return info.Size()
}
