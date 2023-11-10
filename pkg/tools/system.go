package tools

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/goravel/framework/facades"
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
func Read(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		facades.Log().With(map[string]any{
			"path":  path,
			"error": err.Error(),
		}).Tags("面板", "工具函数").Info("读取文件失败")
		return ""
	}

	return string(data)
}

// Remove 删除文件/目录
func Remove(path string) bool {
	if err := os.RemoveAll(path); err != nil {
		facades.Log().With(map[string]any{
			"path":  path,
			"error": err.Error(),
		}).Tags("面板", "工具函数").Info("删除文件/目录失败")
		return false
	}

	return true
}

// Exec 执行 shell 命令
func Exec(shell string) string {
	cmd := exec.Command("bash", "-c", shell)

	output, err := cmd.CombinedOutput()
	if err != nil {
		if support.Env == support.EnvTest {
			fmt.Println(string(output))
			fmt.Println(err.Error())
			panic(err)
		} else {
			facades.Log().With(map[string]any{
				"shell": shell,
				"error": err.Error(),
			}).Tags("面板", "工具函数").Info("执行命令失败")
		}
		return ""
	}

	return strings.TrimSpace(string(output))
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
			} else {
				facades.Log().With(map[string]any{
					"shell": shell,
					"error": err.Error(),
				}).Tags("面板", "工具函数").Info("异步执行命令失败")
			}
		}
	}()

	return nil
}

// Mkdir 创建目录
func Mkdir(path string, permission os.FileMode) bool {
	if err := os.MkdirAll(path, permission); err != nil {
		facades.Log().With(map[string]any{
			"path":       path,
			"permission": permission,
			"error":      err.Error(),
		}).Tags("面板", "工具函数").Info("创建目录失败")
		return false
	}

	return true
}

// Chmod 修改文件/目录权限
func Chmod(path string, permission os.FileMode) bool {
	if err := os.Chmod(path, permission); err != nil {
		facades.Log().With(map[string]any{
			"path":       path,
			"permission": permission,
		}).Tags("面板", "工具函数").Info("修改文件/目录权限失败")
		return false
	}

	return true
}

// Chown 修改文件/目录所有者
func Chown(path, user, group string) bool {
	cmd := exec.Command("chown", "-R", user+":"+group, path)

	err := cmd.Run()
	if err != nil {
		facades.Log().With(map[string]any{
			"path":  path,
			"user":  user,
			"group": group,
			"error": err.Error(),
		}).Tags("面板", "工具函数").Info("修改文件/目录所有者失败")
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

// Mv 移动文件/目录
func Mv(src, dst string) (bool, error) {
	cmd := exec.Command("mv", src, dst)

	err := cmd.Run()
	if err != nil {
		return false, err
	}

	return true, nil
}

// Cp 复制文件/目录
func Cp(src, dst string) (bool, error) {
	cmd := exec.Command("cp", "-r", src, dst)

	err := cmd.Run()
	if err != nil {
		return false, err
	}

	return true, nil
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

	if err != nil {
		return 0, err
	}

	return size, nil
}

// FileSize 获取文件大小
func FileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	return info.Size(), nil
}
