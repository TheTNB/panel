package tools

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
	"strings"
	"text/template"
	"unicode/utf8"
)

// FirstElement 安全地获取 args[0]，避免 panic: runtime error: index out of range
func FirstElement(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

// RandomNumber 生成长度为 length 随机数字字符串
func RandomNumber(length int) string {
	table := [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, b, length)
	if n != length {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

// RandomString 生成长度为 length 的随机字符串
func RandomString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	letters := "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i, v := range b {
		b[i] = letters[v%byte(len(letters))]
	}
	return string(b)
}

// MD5 生成字符串的 MD5 值
func MD5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

// FormatBytes 格式化bytes
func FormatBytes(size float64) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}

	i := 0
	for ; size >= 1024 && i < len(units); i++ {
		size /= 1024
	}

	return fmt.Sprintf("%.2f %s", size, units[i])
}

// Cut 裁剪字符串
func Cut(str, begin, end string) string {
	bIndex := strings.Index(str, begin)
	eIndex := strings.Index(str, end)
	if bIndex == -1 || eIndex == -1 || bIndex > eIndex {
		return ""
	}

	b := utf8.RuneCountInString(str[:bIndex]) + utf8.RuneCountInString(begin)
	e := utf8.RuneCountInString(str[:eIndex])
	if b > e {
		return ""
	}

	return string([]rune(str)[b:e])
}

// Escape 转义字符串
func Escape(str string) string {
	return template.HTMLEscapeString(str)
}
