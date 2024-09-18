package copier

import (
	"encoding/json"
	"fmt"
)

func Copy[T any](from any) (*T, error) {
	to := new(T)
	b, err := json.Marshal(from)
	if err != nil {
		return nil, fmt.Errorf("copier: marshal data err: %w", err)
	}
	if err = json.Unmarshal(b, to); err != nil {
		return nil, fmt.Errorf("copier: unmarshal data err: %w", err)
	}
	return to, nil
}
