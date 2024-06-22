package slice

import "github.com/spf13/cast"

// ToAny 将任意类型切片转换为 []any
func ToAny[T any](s []T) []any {
	result := make([]any, len(s))
	for i, v := range s {
		result[i] = v
	}
	return result
}

// ToString 将任意类型切片转换为 []string
func ToString[T any](s []T) []string {
	result := make([]string, len(s))
	for i, v := range s {
		result[i] = cast.ToString(v)
	}
	return result
}

// ToInt 将任意类型切片转换为 []int
func ToInt[T any](s []T) []int {
	result := make([]int, len(s))
	for i, v := range s {
		result[i] = cast.ToInt(v)
	}
	return result
}
