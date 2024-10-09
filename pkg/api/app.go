package api

import (
	"fmt"
	"time"
)

type App struct {
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Slug        string    `json:"slug"`
	Icon        string    `json:"icon"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Categories  []string  `json:"categories"`
	Depends     string    `json:"depends"`
	Channels    []struct {
		Slug      string `json:"slug"`
		Name      string `json:"name"`
		Panel     string `json:"panel"`
		Install   string `json:"install"`
		Uninstall string `json:"uninstall"`
		Update    string `json:"update"`
		Subs      []struct {
			Log     string `json:"log"`
			Version string `json:"version"`
		} `json:"subs"`
	} `json:"channels"`
	Order int `json:"order"`
}

type Apps []*App

// Apps 返回所有应用
func (r *API) Apps() (*Apps, error) {
	resp, err := r.client.R().SetResult(&Response{}).Get("/apps")
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, fmt.Errorf("failed to get apps: %s", resp.String())
	}

	apps, err := getResponseData[Apps](resp)
	if err != nil {
		return nil, err
	}

	return apps, nil
}

// AppBySlug 根据slug返回应用
func (r *API) AppBySlug(slug string) (*App, error) {
	resp, err := r.client.R().SetResult(&Response{}).Get(fmt.Sprintf("/apps/%s", slug))
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, fmt.Errorf("failed to get app: %s", resp.String())
	}

	app, err := getResponseData[App](resp)
	if err != nil {
		return nil, err
	}

	return app, nil
}
