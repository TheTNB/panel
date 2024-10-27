package types

import "strings"

type NV struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type LV struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type LVInt struct {
	Label string `json:"label"`
	Value int    `json:"value"`
}

// KVToMap 将 key-value 切片转换为 map
func KVToMap(kvs []KV) map[string]string {
	m := make(map[string]string)
	for _, item := range kvs {
		m[item.Key] = item.Value
	}

	return m
}

// MapToKV 将 map 转换为 key-value 切片
func MapToKV(m map[string]string) []KV {
	kvs := make([]KV, 0)
	for k, v := range m {
		kvs = append(kvs, KV{Key: k, Value: v})
	}

	return kvs
}

// KVToSlice 将 key-value 切片转换为 key=value 切片
func KVToSlice(kvs []KV) []string {
	s := make([]string, 0)
	for _, item := range kvs {
		s = append(s, item.Key+"="+item.Value)
	}

	return s
}

// SliceToKV 将 key=value 切片转换为 key-value 切片
func SliceToKV(s []string) []KV {
	kvs := make([]KV, 0)
	for _, item := range s {
		kv := strings.SplitN(item, "=", 2)
		if len(kv) == 2 {
			kvs = append(kvs, KV{Key: kv[0], Value: kv[1]})
		}
	}

	return kvs
}
