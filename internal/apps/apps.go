package apps

import (
	"reflect"
	"slices"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/wire"

	"github.com/TheTNB/panel/internal/apps/benchmark"
	"github.com/TheTNB/panel/internal/apps/docker"
	"github.com/TheTNB/panel/internal/apps/fail2ban"
	"github.com/TheTNB/panel/internal/apps/frp"
	"github.com/TheTNB/panel/internal/apps/gitea"
	"github.com/TheTNB/panel/internal/apps/memcached"
	"github.com/TheTNB/panel/internal/apps/mysql"
	"github.com/TheTNB/panel/internal/apps/nginx"
	"github.com/TheTNB/panel/internal/apps/php"
	"github.com/TheTNB/panel/internal/apps/phpmyadmin"
	"github.com/TheTNB/panel/internal/apps/podman"
	"github.com/TheTNB/panel/internal/apps/postgresql"
	"github.com/TheTNB/panel/internal/apps/pureftpd"
	"github.com/TheTNB/panel/internal/apps/redis"
	"github.com/TheTNB/panel/internal/apps/rsync"
	"github.com/TheTNB/panel/internal/apps/s3fs"
	"github.com/TheTNB/panel/internal/apps/supervisor"
	"github.com/TheTNB/panel/internal/apps/toolbox"
)

var ProviderSet = wire.NewSet(
	NewLoader,
	benchmark.NewApp,
	docker.NewApp,
	fail2ban.NewApp,
	frp.NewApp,
	gitea.NewApp,
	memcached.NewApp,
	mysql.NewApp,
	nginx.NewApp,
	php.NewApp,
	phpmyadmin.NewApp,
	podman.NewApp,
	postgresql.NewApp,
	pureftpd.NewApp,
	redis.NewApp,
	rsync.NewApp,
	s3fs.NewApp,
	supervisor.NewApp,
	toolbox.NewApp,
)

var slugs []string

type Loader struct {
	benchmark  *benchmark.App
	docker     *docker.App
	fail2ban   *fail2ban.App
	frp        *frp.App
	gitea      *gitea.App
	memcached  *memcached.App
	mysql      *mysql.App
	nginx      *nginx.App
	php        *php.App
	phpmyadmin *phpmyadmin.App
	podman     *podman.App
	postgresql *postgresql.App
	pureftpd   *pureftpd.App
	redis      *redis.App
	rsync      *rsync.App
	s3fs       *s3fs.App
	supervisor *supervisor.App
	toolbox    *toolbox.App
}

func NewLoader(
	benchmark *benchmark.App,
	docker *docker.App,
	fail2ban *fail2ban.App,
	frp *frp.App,
	gitea *gitea.App,
	memcached *memcached.App,
	mysql *mysql.App,
	nginx *nginx.App,
	php *php.App,
	phpmyadmin *phpmyadmin.App,
	podman *podman.App,
	postgresql *postgresql.App,
	pureftpd *pureftpd.App,
	redis *redis.App,
	rsync *rsync.App,
	s3fs *s3fs.App,
	supervisor *supervisor.App,
	toolbox *toolbox.App,
) *Loader {
	loader := &Loader{
		benchmark:  benchmark,
		docker:     docker,
		fail2ban:   fail2ban,
		frp:        frp,
		gitea:      gitea,
		memcached:  memcached,
		mysql:      mysql,
		nginx:      nginx,
		php:        php,
		phpmyadmin: phpmyadmin,
		podman:     podman,
		postgresql: postgresql,
		pureftpd:   pureftpd,
		redis:      redis,
		rsync:      rsync,
		s3fs:       s3fs,
		supervisor: supervisor,
		toolbox:    toolbox,
	}

	loader.initSlugs()
	return loader
}

func (r *Loader) Register(mux chi.Router) {
	mux.Route("/benchmark", r.benchmark.Route)
	mux.Route("/docker", r.docker.Route)
	mux.Route("/fail2ban", r.fail2ban.Route)
	mux.Route("/frp", r.frp.Route)
	mux.Route("/gitea", r.gitea.Route)
	mux.Route("/memcached", r.memcached.Route)
	mux.Route("/mysql", r.mysql.Route)
	mux.Route("/nginx", r.nginx.Route)
	mux.Route("/php74", r.php.Route(74))
	mux.Route("/php80", r.php.Route(80))
	mux.Route("/php81", r.php.Route(81))
	mux.Route("/php82", r.php.Route(82))
	mux.Route("/php83", r.php.Route(83))
	mux.Route("/php84", r.php.Route(84))
	mux.Route("/phpmyadmin", r.phpmyadmin.Route)
	mux.Route("/podman", r.podman.Route)
	mux.Route("/postgresql", r.postgresql.Route)
	mux.Route("/pureftpd", r.pureftpd.Route)
	mux.Route("/redis", r.redis.Route)
	mux.Route("/rsync", r.rsync.Route)
	mux.Route("/s3fs", r.s3fs.Route)
	mux.Route("/supervisor", r.supervisor.Route)
	mux.Route("/toolbox", r.toolbox.Route)
}

func (r *Loader) initSlugs() []string {
	if len(slugs) == 0 {
		v := reflect.Indirect(reflect.ValueOf(r))
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			slug := strings.ToLower(v.Type().Field(i).Name)
			if !field.IsNil() {
				slugs = append(slugs, slug)
			}
		}

		// 处理php
		slugs = slices.DeleteFunc(slugs, func(slug string) bool {
			return slug == "php"
		})
		slugs = append(slugs, "php74", "php80", "php81", "php82", "php83", "php84")
	}

	return slugs
}

func Slugs() []string {
	return slugs
}
