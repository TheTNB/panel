package io

import (
	"errors"
	"github.com/go-rat/utils/debug"
	"os/exec"
	"path/filepath"
	"strings"
)

type FormatArchive string

const (
	Zip      FormatArchive = "zip"
	Gz       FormatArchive = "gz"
	Bz2      FormatArchive = "bz2"
	Tar      FormatArchive = "tar"
	TarGz    FormatArchive = "tar.gz"
	Xz       FormatArchive = "xz"
	SevenZip FormatArchive = "7z"
)

// Compress 压缩文件
func Compress(dir string, src []string, dst string) error {
	if len(src) == 0 {
		src = append(src, ".")
	}
	// 去掉路径前缀，减少压缩包内文件夹层级
	for i, s := range src {
		src[i] = strings.TrimPrefix(s, dir)
		if src[i] == "" {
			src[i] = "."
		}
	}

	cmd := new(exec.Cmd)
	cmd.Dir = dir

	format, err := formatArchiveByPath(dst)
	if err != nil {
		return err
	}

	switch format {
	case Zip:
		cmd = exec.Command("zip", append([]string{"-qr", "-o", dst}, src...)...)
	case Gz:
		cmd = exec.Command("tar", append([]string{"-czf", dst}, src...)...)
	case Bz2:
		cmd = exec.Command("tar", append([]string{"-cjf", dst}, src...)...)
	case Tar:
		cmd = exec.Command("tar", append([]string{"-cf", dst}, src...)...)
	case TarGz:
		cmd = exec.Command("tar", append([]string{"-czf", dst}, src...)...)
	case Xz:
		cmd = exec.Command("tar", append([]string{"-cJf", dst}, src...)...)
	case SevenZip:
		cmd = exec.Command("7z", append([]string{"a", "-y", dst}, src...)...)
	default:
		return errors.New("unsupported format")
	}

	out, err := cmd.CombinedOutput()
	debug.Dump(string(out))

	return err
}

// UnCompress 解压文件
func UnCompress(src string, dst string) error {
	var cmd *exec.Cmd

	format, err := formatArchiveByPath(src)
	if err != nil {
		return err
	}

	switch format {
	case Zip:
		cmd = exec.Command("unzip", "-qo", src, "-d", dst)
	case Gz:
		cmd = exec.Command("tar", "-xzf", src, "-C", dst)
	case Bz2:
		cmd = exec.Command("tar", "-xjf", src, "-C", dst)
	case Tar:
		cmd = exec.Command("tar", "-xf", src, "-C", dst)
	case TarGz:
		cmd = exec.Command("tar", "-xzf", src, "-C", dst)
	case Xz:
		cmd = exec.Command("tar", "-xJf", src, "-C", dst)
	case SevenZip:
		cmd = exec.Command("7z", "x", "-y", src, "-o"+dst)
	default:
		return errors.New("unsupported format")
	}

	return cmd.Run()
}

// formatArchiveByPath 根据文件后缀获取压缩格式
func formatArchiveByPath(path string) (FormatArchive, error) {
	switch filepath.Ext(path) {
	case ".zip":
		return Zip, nil
	case ".gz":
		return Gz, nil
	case ".bz2":
		return Bz2, nil
	case ".tar":
		return Tar, nil
	case ".tar.gz":
		return TarGz, nil
	case ".xz":
		return Xz, nil
	case ".7z":
		return SevenZip, nil
	default:
		return "", errors.New("unknown format")
	}
}
