package shell

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/goravel/framework/support"

	"github.com/TheTNB/panel/v2/pkg/slice"
)

// Execf 执行 shell 命令
func Execf(shell string, args ...any) (string, error) {
	if !CheckArgs(slice.ToString(args)...) {
		return "", errors.New("发现危险的命令参数，中止执行")
	}

	var cmd *exec.Cmd
	_ = os.Setenv("LC_ALL", "C")
	cmd = exec.Command("bash", "-c", fmt.Sprintf(shell, args...))

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", errors.New(strings.TrimSpace(stderr.String()))
	}

	return strings.TrimSpace(stdout.String()), err
}

// ExecfAsync 异步执行 shell 命令
func ExecfAsync(shell string, args ...any) error {
	if !CheckArgs(slice.ToString(args)...) {
		return errors.New("发现危险的命令参数，中止执行")
	}

	var cmd *exec.Cmd
	_ = os.Setenv("LC_ALL", "C")
	cmd = exec.Command("bash", "-c", fmt.Sprintf(shell, args...))

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

// ExecfWithTimeout 执行 shell 命令并设置超时时间
func ExecfWithTimeout(timeout time.Duration, shell string, args ...any) (string, error) {
	if !CheckArgs(slice.ToString(args)...) {
		return "", errors.New("发现危险的命令参数，中止执行")
	}

	var cmd *exec.Cmd
	_ = os.Setenv("LC_ALL", "C")
	cmd = exec.Command("bash", "-c", fmt.Sprintf(shell, args...))

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Start()
	if err != nil {
		return "", err
	}

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(timeout):
		_ = cmd.Process.Kill()
		return "", errors.New("执行超时")
	case err = <-done:
		if err != nil {
			return "", errors.New(strings.TrimSpace(stderr.String()))
		}
	}

	return strings.TrimSpace(stdout.String()), err
}

// CheckArgs 检查危险的参数
func CheckArgs(args ...string) bool {
	if len(args) == 0 {
		return true
	}

	dangerous := []string{"&", "|", ";", "$", "'", `"`, "(", ")", "`", "\n", "\r", ">", "<", "{", "}", "[", "]", "\\"}
	for _, arg := range args {
		for _, char := range dangerous {
			if strings.Contains(arg, char) {
				return false
			}
		}
	}

	return true
}
