package data

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

func getDockerClient(sock string) *resty.Client {
	client := resty.New()
	client.SetTimeout(1 * time.Minute)
	client.SetRetryCount(2)
	client.SetTransport(&http.Transport{
		DialContext: func(ctx context.Context, _ string, _ string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "unix", sock)
		},
	})
	client.SetBaseURL("http://d/v1.40")
	return client
}
