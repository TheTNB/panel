package shell

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Execf 执行 shell 命令
func Execf(shell string, args ...any) (string, error) {
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
	var cmd *exec.Cmd
	_ = os.Setenv("LC_ALL", "C")
	cmd = exec.Command("bash", "-c", fmt.Sprintf(shell, args...))

	err := cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		if err = cmd.Wait(); err != nil {
			fmt.Println(err.Error())
		}
	}()

	return nil
}

// ExecfWithTimeout 执行 shell 命令并设置超时时间
func ExecfWithTimeout(timeout time.Duration, shell string, args ...any) (string, error) {
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

// ExecfWithOutput 执行 shell 命令并输出到终端
func ExecfWithOutput(shell string, args ...any) error {
	_ = os.Setenv("LC_ALL", "C")
	cmd := exec.Command("bash", "-c", fmt.Sprintf(shell, args...))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
