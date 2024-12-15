package app

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/bddjr/hlfhr"
	"github.com/go-chi/chi/v5"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/gookit/validate"
	"github.com/knadh/koanf/v2"
	"github.com/robfig/cron/v3"
)

type Web struct {
	conf     *koanf.Koanf
	router   *chi.Mux
	server   *hlfhr.Server
	migrator *gormigrate.Gormigrate
	cron     *cron.Cron
}

func NewWeb(conf *koanf.Koanf, router *chi.Mux, server *hlfhr.Server, migrator *gormigrate.Gormigrate, cron *cron.Cron, _ *validate.Validation) *Web {
	return &Web{
		conf:     conf,
		router:   router,
		server:   server,
		migrator: migrator,
		cron:     cron,
	}
}

func (r *Web) Run() error {
	// migrate database
	if err := r.migrator.Migrate(); err != nil {
		return err
	}
	fmt.Println("[DB] database migrated")

	// start cron scheduler
	r.cron.Start()
	fmt.Println("[CRON] cron scheduler started")

	// run http server
	if r.conf.Bool("http.tls") {
		cert := filepath.Join(Root, "panel/storage/cert.pem")
		key := filepath.Join(Root, "panel/storage/cert.key")
		fmt.Println("[HTTP] listening and serving on port", r.conf.MustInt("http.port"), "with tls")
		if err := r.server.ListenAndServeTLS(cert, key); !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	} else {
		fmt.Println("[HTTP] listening and serving on port", r.conf.MustInt("http.port"))
		if err := r.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}

	return nil
}
