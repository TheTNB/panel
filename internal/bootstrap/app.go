package bootstrap

import (
	"os"
)

func boot() {
	if os.Geteuid() != 0 {
		panic("panel must run as root")
	}

	initConf()
	initGlobal()
	initLogger()
	initOrm()
	runMigrate()
	bootCrypter()
}

func BootWeb() {
	boot()
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
