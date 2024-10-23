package bootstrap

import (
	"os"
	"runtime/debug"
)

func boot() {
	if os.Geteuid() != 0 {
		panic("panel must run as root")
	}

	debug.SetGCPercent(10)
	debug.SetMemoryLimit(64 << 20)

	initConf()
	initGlobal()
	initLogger()
	initOrm()
	runMigrate()
}

func BootWeb() {
	boot()
	initValidator()
	initSession()
	initQueue()
	initCron()
	initHttp()

	select {}
}

func BootCli() {
	boot()
	initCli()
}
