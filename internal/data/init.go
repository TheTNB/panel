package data

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/go-rat/utils/hash"
	"github.com/go-resty/resty/v2"
	"github.com/samber/do/v2"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/pkg/api"
	"github.com/TheTNB/panel/pkg/os"
)

var injector = do.New()

func init() {
	do.Provide(injector, func(i do.Injector) (biz.AppRepo, error) {
		return &appRepo{
			api: api.NewAPI(app.Version),
		}, nil
	})
	do.Provide(injector, func(i do.Injector) (biz.BackupRepo, error) {
		return &backupRepo{}, nil
	})
	do.Provide(injector, func(i do.Injector) (biz.CacheRepo, error) {
		return &cacheRepo{}, nil
	})
	do.Provide(injector, func(i do.Injector) (biz.CertRepo, error) {
		return &certRepo{}, nil
	})
	do.Provide(injector, func(i do.Injector) (biz.CertAccountRepo, error) {
		return &certAccountRepo{}, nil
	})
	do.Provide(injector, func(i do.Injector) (biz.CertDNSRepo, error) {
		return &certDNSRepo{}, nil
	})
	do.Provide(injector, func(i do.Injector) (biz.ContainerRepo, error) {
		return &containerRepo{
			client: getDockerClient("/var/run/docker.sock"),
		}, nil
	})
	do.Provide(injector, func(i do.Injector) (biz.ContainerImageRepo, error) {
		return &containerImageRepo{
			client: getDockerClient("/var/run/docker.sock"),
		}, nil
	})
	do.Provide(injector, func(i do.Injector) (biz.ContainerNetworkRepo, error) {
		return &containerNetworkRepo{
			client: getDockerClient("/var/run/docker.sock"),
		}, nil
	})
	do.Provide(injector, func(i do.Injector) (biz.ContainerVolumeRepo, error) {
		return &containerVolumeRepo{
			client: getDockerClient("/var/run/docker.sock"),
		}, nil
	})
	do.Provide(injector, func(i do.Injector) (biz.CronRepo, error) {
		return &cronRepo{}, nil
	})
	do.Provide(injector, func(i do.Injector) (biz.DatabaseServerRepo, error) {
		return &databaseServerRepo{}, nil
	})
	do.Provide(injector, func(i do.Injector) (biz.DatabaseUserRepo, error) {
		return &databaseUserRepo{}, nil
	})
	do.Provide(injector, func(i do.Injector) (biz.DatabaseRepo, error) {
		return &databaseRepo{}, nil
	})
	do.Provide(injector, func(i do.Injector) (biz.MonitorRepo, error) {
		return &monitorRepo{}, nil
	})
	do.Provide(injector, func(i do.Injector) (biz.SafeRepo, error) {
		var ssh string
		if os.IsRHEL() {
			ssh = "sshd"
		} else {
			ssh = "ssh"
		}
		return &safeRepo{
			ssh: ssh,
		}, nil
	})
	do.Provide(injector, func(i do.Injector) (biz.SettingRepo, error) {
		return &settingRepo{}, nil
	})
	do.Provide(injector, func(i do.Injector) (biz.SSHRepo, error) {
		return &sshRepo{}, nil
	})
	do.Provide(injector, func(i do.Injector) (biz.TaskRepo, error) {
		task := &taskRepo{}
		task.DispatchWaiting()
		return task, nil
	})
	do.Provide(injector, func(i do.Injector) (biz.UserRepo, error) {
		return &userRepo{
			hasher: hash.NewArgon2id(),
		}, nil
	})
	do.Provide(injector, func(i do.Injector) (biz.WebsiteRepo, error) {
		return &websiteRepo{}, nil
	})
}

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
