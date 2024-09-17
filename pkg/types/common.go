package types

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

// KVToMap 将 key-value 切片转换为 map
func KVToMap(kvs []KV) map[string]string {
	m := make(map[string]string)
	for _, item := range kvs {
		m[item.Key] = item.Value
	}

	return m
}

// KVToSlice 将 key-value 切片转换为 key=value 切片
func KVToSlice(kvs []KV) []string {
	var s []string
	for _, item := range kvs {
		s = append(s, item.Key+"="+item.Value)
	}

	return s
}
