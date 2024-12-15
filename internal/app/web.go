package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/bddjr/hlfhr"
	"github.com/cloudflare/tableflip"
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
	if runtime.GOOS != "windows" {
		return r.runServer()
	}

	return r.runServerFallback()
}

// runServer graceful run server
func (r *Web) runServer() error {
	upg, err := tableflip.New(tableflip.Options{})
	if err != nil {
		return err
	}
	defer upg.Stop()

	// By prefixing PID to log, easy to interrupt from another process.
	log.SetPrefix(fmt.Sprintf("[PID %d]", os.Getpid()))

	// Listen for the process signal to trigger the tableflip upgrade.
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP)
		for range sig {
			if err = upg.Upgrade(); err != nil {
				log.Println("[Graceful] upgrade failed:", err)
			}
		}
	}()

	ln, err := upg.Listen("tcp", r.conf.MustString("http.address"))
	if err != nil {
		return err
	}
	defer ln.Close()

	fmt.Println("[HTTP] listening and serving on", r.conf.MustString("http.address"))
	go func() {
		if err = r.server.Serve(ln); !errors.Is(err, http.ErrServerClosed) {
			log.Println("[HTTP] server error:", err)
		}
	}()

	// tableflip ready
	if err = upg.Ready(); err != nil {
		return err
	}

	fmt.Println("[Graceful] ready for upgrade")
	<-upg.Exit()

	// Make sure to set a deadline on exiting the process
	// after upg.Exit() is closed. No new upgrades can be
	// performed if the parent doesn't exit.
	time.AfterFunc(60*time.Second, func() {
		log.Println("[Graceful] shutdown timeout, force exit")
		os.Exit(1)
	})

	// Wait for connections to drain.
	return r.server.Shutdown(context.Background())
}

// runServerFallback fallback for windows
func (r *Web) runServerFallback() error {
	fmt.Println("[HTTP] listening and serving on", r.conf.MustString("http.address"))
	if err := r.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
