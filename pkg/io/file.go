package io

import (
	"archive/zip"
	"context"
	"errors"
	stdio "io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver/v4"
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

var (
	ErrFormatNotSupport = errors.New("不支持此格式")
	ErrNotDirectory     = errors.New("目标不是目录")
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

// FileInfo 获取文件大小
func FileInfo(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

// Compress 压缩文件
func Compress(src []string, dst string, format FormatArchive) error {
	// 不支持7z
	if format == SevenZip {
		return ErrFormatNotSupport
	}
	arch := getFormat(format)

	srcMap := make(map[string]string, len(src))
	for _, s := range src {
		base := filepath.Base(s)
		srcMap[s] = base
	}

	dir := filepath.Dir(dst)
	if !Exists(dir) {
		if err := Mkdir(dir, 0755); err != nil {
			return err
		}
	}

	files, err := archiver.FilesFromDisk(nil, srcMap)
	if err != nil {
		return err
	}

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer out.Close()

	err = arch.Archive(context.Background(), out, files)
	if err != nil {
		_ = Remove(dst)
	}

	return nil
}

// UnCompress 解压文件
func UnCompress(src string, dst string, format FormatArchive) error {
	handler := func(ctx context.Context, f archiver.File) error {
		info := f.FileInfo
		fileName := f.NameInArchive
		filePath := filepath.Join(dst, fileName)

		if f.FileInfo.IsDir() {
			if err := Mkdir(filePath, info.Mode()); err != nil {
				return err
			}
			return nil
		}

		parentDir := path.Dir(filePath)
		if !Exists(parentDir) {
			if err := Mkdir(parentDir, info.Mode()); err != nil {
				return err
			}
		}

		r, err := f.Open()
		if err != nil {
			return err
		}
		w, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, info.Mode())
		if err != nil {
			return err
		}
		defer r.Close()
		defer w.Close()

		if _, err = stdio.Copy(w, r); err != nil {
			return err
		}

		return nil
	}

	arch := getFormat(format)
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	if !Exists(dst) {
		if err = Mkdir(dst, 0755); err != nil {
			return err
		}
	}
	if !IsDir(dst) {
		return ErrNotDirectory
	}

	return arch.Extract(context.Background(), file, nil, handler)
}

// FormatArchiveByPath 根据文件后缀获取压缩格式
func FormatArchiveByPath(path string) (FormatArchive, error) {
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

// TempFile 创建临时文件
func TempFile(dir, prefix string) (*os.File, error) {
	return os.CreateTemp(dir, prefix)
}

func getFormat(f FormatArchive) archiver.CompressedArchive {
	format := archiver.CompressedArchive{}
	switch f {
	case Tar:
		format.Archival = archiver.Tar{}
	case TarGz, Gz:
		format.Compression = archiver.Gz{}
		format.Archival = archiver.Tar{}
	case Zip:
		format.Archival = archiver.Zip{
			Compression: zip.Deflate,
		}
	case Bz2:
		format.Compression = archiver.Bz2{}
		format.Archival = archiver.Tar{}
	case Xz:
		format.Compression = archiver.Xz{}
		format.Archival = archiver.Tar{}
	case SevenZip:
		format.Archival = archiver.SevenZip{}

	}
	return format
}
