package shell

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"
	"strings"
	"time"
)

// Execf 执行 shell 命令
func Execf(shell string, args ...any) (string, error) {
	if !preCheckArg(args) {
		return "", errors.New("command contains illegal characters")
	}

	_ = os.Setenv("LC_ALL", "C")
	cmd := exec.Command("bash", "-c", fmt.Sprintf(shell, args...))

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return strings.TrimSpace(stdout.String()), errors.New(strings.TrimSpace(stderr.String()))
	}

	return strings.TrimSpace(stdout.String()), err
}

// ExecfAsync 异步执行 shell 命令
func ExecfAsync(shell string, args ...any) error {
	if !preCheckArg(args) {
		return errors.New("command contains illegal characters")
	}

	_ = os.Setenv("LC_ALL", "C")
	cmd := exec.Command("bash", "-c", fmt.Sprintf(shell, args...))

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
	if !preCheckArg(args) {
		return "", errors.New("command contains illegal characters")
	}

	_ = os.Setenv("LC_ALL", "C")
	cmd := exec.Command("bash", "-c", fmt.Sprintf(shell, args...))

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
		return strings.TrimSpace(stdout.String()), errors.New("执行超时")
	case err = <-done:
		if err != nil {
			return strings.TrimSpace(stdout.String()), errors.New(strings.TrimSpace(stderr.String()))
		}
	}

	return strings.TrimSpace(stdout.String()), err
}

// ExecfWithOutput 执行 shell 命令并输出到终端
func ExecfWithOutput(shell string, args ...any) error {
	if !preCheckArg(args) {
		return errors.New("command contains illegal characters")
	}

	_ = os.Setenv("LC_ALL", "C")
	cmd := exec.Command("bash", "-c", fmt.Sprintf(shell, args...))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// ExecfWithPipe 执行 shell 命令并返回管道
func ExecfWithPipe(ctx context.Context, shell string, args ...any) (out io.ReadCloser, err error) {
	if !preCheckArg(args) {
		return nil, errors.New("command contains illegal characters")
	}

	_ = os.Setenv("LC_ALL", "C")
	cmd := exec.CommandContext(ctx, "bash", "-c", fmt.Sprintf(shell, args...))

	out, err = cmd.StdoutPipe()
	if err != nil {
		return
	}

	cmd.Stderr = cmd.Stdout
	err = cmd.Start()
	return
}

func preCheckArg(args []any) bool {
	illegals := []any{`&`, `|`, `;`, `$`, `'`, `"`, "`", `(`, `)`, "\n", "\r", `>`, `<`}
	for arg := range slices.Values(args) {
		if slices.Contains(illegals, arg) {
			return false
		}
	}

	return true
}
