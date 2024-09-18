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
	Requires    []string  `json:"requires"`
	Excludes    []string  `json:"excludes"`
	Versions    []struct {
		Version      string `json:"version"`
		Install      string `json:"install"`
		Uninstall    string `json:"uninstall"`
		Update       string `json:"update"`
		PanelVersion string `json:"panel_version"`
	} `json:"versions"`
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
