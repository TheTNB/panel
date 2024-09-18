package api

import (
	"fmt"
	"time"
)

type Rewrite struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
}

type Rewrites []Rewrite

func (r *API) RewritesByType(typ string) (*Rewrites, error) {
	resp, err := r.client.R().SetResult(&Response{}).Get(fmt.Sprintf("/rewrites/%s", typ))
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, fmt.Errorf("failed to get rewrites: %s", resp.String())
	}

	rewrites, err := getResponseData[Rewrites](resp)
	if err != nil {
		return nil, err
	}

	return rewrites, nil
}
