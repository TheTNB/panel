package tools

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/goravel/framework/support"
	"github.com/goravel/framework/support/env"
	"github.com/mholt/archiver/v3"
	"github.com/spf13/cast"
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

// Remove 删除文件/目录
func Remove(path string) error {
	return os.RemoveAll(path)
}

// Exec 执行 shell 命令
func Exec(shell string) (string, error) {
	var cmd *exec.Cmd
	if env.IsLinux() {
		cmd = exec.Command("bash", "-c", "LC_ALL=C "+shell)
	} else {
		cmd = exec.Command("cmd", "/C", "chcp 65001 >nul && "+shell)
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()
	if err != nil {
		return "", errors.New(strings.TrimSpace(stderrBuf.String()))
	}

	return strings.TrimSpace(stdoutBuf.String()), err
}

// ExecAsync 异步执行 shell 命令
func ExecAsync(shell string) error {
	var cmd *exec.Cmd
	if env.IsLinux() {
		cmd = exec.Command("bash", "-c", "LC_ALL=C "+shell)
	} else {
		cmd = exec.Command("cmd", "/C", "chcp 65001 >nul && "+shell)
	}

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
func Chmod(path string, permission uint) error {
	if env.IsWindows() {
		return errors.New("chmod is not supported on Windows")
	}

	cmd := exec.Command("chmod", "-R", cast.ToString(permission), path)
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

// TempDir 创建临时目录
func TempDir(prefix string) (string, error) {
	return os.MkdirTemp("", prefix)
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

// GetUser 通过 uid 获取用户名
func GetUser(uid uint32) string {
	usr, err := user.LookupId(cast.ToString(uid))
	if err != nil {
		return ""
	}
	return usr.Username
}

// GetGroup 通过 gid 获取组名
func GetGroup(gid uint32) string {
	usr, err := user.LookupGroupId(cast.ToString(gid))
	if err != nil {
		return ""
	}
	return usr.Name
}
