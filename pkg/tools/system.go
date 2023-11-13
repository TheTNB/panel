package tools

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/goravel/framework/support"
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

// Read 读取文件
// TODO 重构带 error 返回
func Read(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	return string(data)
}

// Remove 删除文件/目录
// TODO 重构带 error 返回
func Remove(path string) bool {
	if err := os.RemoveAll(path); err != nil {
		return false
	}

	return true
}

// Exec 执行 shell 命令
func Exec(shell string) (string, error) {
	cmd := exec.Command("bash", "-c", "LC_ALL=C "+shell)

	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

// ExecAsync 异步执行 shell 命令
func ExecAsync(shell string) error {
	cmd := exec.Command("bash", "-c", shell)
	err := cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		err := cmd.Wait()
		if err != nil {
			if support.Env == support.EnvTest {
				fmt.Println(err.Error())
				panic(err)
			}
		}
	}()

	return nil
}

// Mkdir 创建目录
func Mkdir(path string, permission os.FileMode) error {
	return os.MkdirAll(path, permission)
}

// Chmod 修改文件/目录权限
func Chmod(path string, permission os.FileMode) error {
	return os.Chmod(path, permission)
}

// Chown 修改文件/目录所有者
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

// Mv 移动文件/目录
func Mv(src, dst string) error {
	cmd := exec.Command("mv", src, dst)

	return cmd.Run()
}

// Cp 复制文件/目录
func Cp(src, dst string) error {
	cmd := exec.Command("cp", "-r", src, dst)

	return cmd.Run()
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

// FileInfo 获取文件大小
func FileInfo(path string) (os.FileInfo, error) {
	return os.Stat(path)
}
