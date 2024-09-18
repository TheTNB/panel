package api

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/shirou/gopsutil/host"

	"github.com/TheTNB/panel/pkg/copier"
)

type API struct {
	panelVersion string
	client       *resty.Client
}

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func NewAPI(panelVersion string, url ...string) *API {
	if len(panelVersion) == 0 {
		panic("panel version is required")
	}
	if len(url) == 0 {
		url = append(url, "https://panel.haozi.net/api")
	}

	hostInfo, err := host.Info()
	if err != nil {
		panic(fmt.Sprintf("failed to get host info: %v", err))
	}

	client := resty.New()
	client.SetTimeout(10 * time.Second)
	client.SetBaseURL(url[0])
	client.SetHeader("User-Agent", fmt.Sprintf("rat-panel/%s %s/%s", panelVersion, hostInfo.Platform, hostInfo.PlatformVersion))

	return &API{
		panelVersion: panelVersion,
		client:       client,
	}
}

func getResponseData[T any](resp *resty.Response) (*T, error) {
	raw, ok := resp.Result().(*Response)
	if !ok {
		return nil, fmt.Errorf("failed to get response data: %s", resp.String())
	}

	res, err := copier.Copy[T](raw.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to copy response data: %w", err)
	}

	return res, nil
}
