package helpers

import (
	"os"
	"os/exec"
	"path/filepath"

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
		facades.Log().Errorf("[面板][Helpers] 执行命令 $s 失败: %s", shell, err.Error())
		return ""
	}

	return string(output)
}

// ExecShellAsync 异步执行 shell 命令
func ExecShellAsync(shell string) {
	cmd := exec.Command("bash", "-c", shell)

	err := cmd.Start()
	if err != nil {
		facades.Log().Errorf("[面板][Helpers] 执行命令 $s 失败: %s", shell, err.Error())
	}

	go func() {
		err := cmd.Wait()
		if err != nil {
			facades.Log().Errorf("[面板][Helpers] 执行命令 $s 失败: %s", shell, err.Error())
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
	cmd := exec.Command("sudo chown", "-R", user+":"+group, path)

	err := cmd.Run()
	if err != nil {
		facades.Log().Errorf("[面板][Helpers] 修改路径 %s 所有者失败: %s", path, err.Error())
		return false
	}

	return true
}
