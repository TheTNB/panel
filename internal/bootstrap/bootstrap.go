package bootstrap

import "github.com/google/wire"

// ProviderSet is bootstrap providers.
var ProviderSet = wire.NewSet(NewConf, NewLog, NewCli, NewValidator, NewRouter, NewHttp, NewDB, NewMigrate, NewLoader, NewSession, NewCron, NewQueue)
