package bootstrap

import (
	"runtime/debug"
)

func Boot() {
	debug.SetGCPercent(10)
	debug.SetMemoryLimit(64 << 20)

	initConf()
	initGlobal()
	initOrm()
	runMigrate()
	initValidator()
	initSession()
	initQueue()
	go initHttp()

	select {}
}
