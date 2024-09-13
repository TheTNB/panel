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
