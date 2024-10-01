package bootstrap

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/TheTNB/panel/internal/app"
)

func initLogger() {
	writeSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "storage/logs/app.log",
		MaxSize:    10,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   true,
	})

	level := zapcore.InfoLevel
	if app.Conf.Bool("app.debug") {
		level = zapcore.DebugLevel
	}

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		writeSyncer,
		level,
	)

	logger := zap.New(core)
	defer logger.Sync()
	app.Logger = logger
}
