package io

import (
	"errors"
	"os/exec"
	"path/filepath"
	"strings"
)

type FormatArchive string

const (
	Zip      FormatArchive = "zip"
	Bz2      FormatArchive = "bz2"
	Tar      FormatArchive = "tar"
	TGz      FormatArchive = "tgz"
	Xz       FormatArchive = "xz"
	SevenZip FormatArchive = "7z"
)

// Compress 压缩文件
func Compress(dir string, src []string, dst string) error {
	if !filepath.IsAbs(dir) || !filepath.IsAbs(dst) {
		return errors.New("dir and dst must be absolute path")
	}
	if len(src) == 0 {
		src = append(src, ".")
	}
	// 去掉路径前缀，减少压缩包内文件夹层级
	for i, s := range src {
		if strings.HasPrefix(s, dir) {
			s = strings.TrimPrefix(s, dir)
			src[i] = strings.TrimPrefix(s, "/")
		}
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
	case TGz:
		cmd = exec.Command("tar", append([]string{"-czf", dst}, src...)...)
	case Bz2:
		cmd = exec.Command("tar", append([]string{"-cjf", dst}, src...)...)
	case Tar:
		cmd = exec.Command("tar", append([]string{"-cf", dst}, src...)...)
	case Xz:
		cmd = exec.Command("tar", append([]string{"-cJf", dst}, src...)...)
	case SevenZip:
		cmd = exec.Command("7z", append([]string{"a", "-y", dst}, src...)...)
	default:
		return errors.New("unsupported format")
	}

	return cmd.Run()
}

// UnCompress 解压文件
func UnCompress(src string, dst string) error {
	if !filepath.IsAbs(src) || !filepath.IsAbs(dst) {
		return errors.New("src and dst must be absolute path")
	}
	if !Exists(dst) {
		if err := Mkdir(dst, 0755); err != nil {
			return err
		}
	}

	format, err := formatArchiveByPath(src)
	if err != nil {
		return err
	}

	var cmd *exec.Cmd
	switch format {
	case Zip:
		cmd = exec.Command("unzip", "-qo", src, "-d", dst)
	case TGz:
		cmd = exec.Command("tar", "-xzf", src, "-C", dst)
	case Bz2:
		cmd = exec.Command("tar", "-xjf", src, "-C", dst)
	case Tar:
		cmd = exec.Command("tar", "-xf", src, "-C", dst)
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
	case ".bz2":
		return Bz2, nil
	case ".tar":
		return Tar, nil
	case ".gz", ".tar.gz", ".tgz":
		return TGz, nil
	case ".xz":
		return Xz, nil
	case ".7z":
		return SevenZip, nil
	default:
		return "", errors.New("unknown format")
	}
}
