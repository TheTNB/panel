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
	"syscall"
	"time"

	"github.com/creack/pty"
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

	if err := cmd.Run(); err != nil {
		return strings.TrimSpace(stdout.String()), fmt.Errorf("run %s failed, err: %s", fmt.Sprintf(shell, args...), strings.TrimSpace(stderr.String()))
	}

	return strings.TrimSpace(stdout.String()), nil
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
			fmt.Println(fmt.Errorf("run %s failed, err: %s", fmt.Sprintf(shell, args...), strings.TrimSpace(err.Error())))
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
		return strings.TrimSpace(stdout.String()), fmt.Errorf("run %s failed, err: %s", fmt.Sprintf(shell, args...), strings.TrimSpace(stderr.String()))
	}

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(timeout):
		_ = cmd.Process.Kill()
		return strings.TrimSpace(stdout.String()), fmt.Errorf("run %s failed, err: %s", fmt.Sprintf(shell, args...), "timeout")
	case err = <-done:
		if err != nil {
			return strings.TrimSpace(stdout.String()), fmt.Errorf("run %s failed, err: %s", fmt.Sprintf(shell, args...), strings.TrimSpace(stderr.String()))
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

// ExecfWithDir 在指定目录下执行 shell 命令
func ExecfWithDir(dir, shell string, args ...any) (string, error) {
	if !preCheckArg(args) {
		return "", errors.New("command contains illegal characters")
	}

	_ = os.Setenv("LC_ALL", "C")
	cmd := exec.Command("bash", "-c", fmt.Sprintf(shell, args...))
	cmd.Dir = dir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return strings.TrimSpace(stdout.String()), fmt.Errorf("run %s failed, err: %s", fmt.Sprintf(shell, args...), strings.TrimSpace(stderr.String()))
	}

	return strings.TrimSpace(stdout.String()), nil
}

// ExecfWithTTY 在伪终端下执行 shell 命令
func ExecfWithTTY(shell string, args ...any) (string, error) {
	if !preCheckArg(args) {
		return "", errors.New("command contains illegal characters")
	}

	_ = os.Setenv("LC_ALL", "C")
	cmd := exec.Command("bash", "-i", "-c", fmt.Sprintf(shell, args...))

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stderr = &stderr // https://github.com/creack/pty/issues/147 取 stderr

	f, err := pty.Start(cmd)
	if err != nil {
		return "", fmt.Errorf("run %s failed", fmt.Sprintf(shell, args...))
	}
	defer f.Close()

	if _, err = io.Copy(&out, f); ptyError(err) != nil {
		return "", fmt.Errorf("run %s failed, out: %s, err: %w", fmt.Sprintf(shell, args...), strings.TrimSpace(out.String()), err)
	}
	if stderr.Len() > 0 {
		return "", fmt.Errorf("run %s failed, out: %s", fmt.Sprintf(shell, args...), strings.TrimSpace(stderr.String()))
	}

	return strings.TrimSpace(out.String()), nil
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

// Linux kernel return EIO when attempting to read from a master pseudo
// terminal which no longer has an open slave. So ignore error here.
// See https://github.com/creack/pty/issues/21
func ptyError(err error) error {
	var pathErr *os.PathError
	if !errors.As(err, &pathErr) || !errors.Is(pathErr.Err, syscall.EIO) {
		return err
	}

	return nil
}
