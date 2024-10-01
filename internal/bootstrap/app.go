package bootstrap

import (
	"runtime/debug"
)

func boot() {
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
	go initHttp()

	select {}
}

func BootCli() {
	boot()
	initCli()
}
