package bootstrap

import (
	"log/slog"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/TheTNB/panel/internal/app"
)

func initLogger() {
	writeSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename: filepath.Join(app.Root, "panel/storage/logs/app.log"),
		MaxSize:  10,
		MaxAge:   30,
		Compress: true,
	})

	level := zap.InfoLevel
	if app.Conf.Bool("app.debug") {
		level = zap.DebugLevel
	}

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		writeSyncer,
		level,
	)

	logger := zap.New(core)
	zap.ReplaceGlobals(logger)
	app.Logger = slog.New(zapslog.NewHandler(logger.Core()))
	slog.SetDefault(app.Logger)
}
