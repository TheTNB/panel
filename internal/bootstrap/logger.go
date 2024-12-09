package bootstrap

import (
	"log/slog"
	"path/filepath"

	"github.com/TheTNB/panel/internal/app"
	"gopkg.in/natefinch/lumberjack.v2"
)

func initLogger() {
	ljLogger := &lumberjack.Logger{
		Filename: filepath.Join(app.Root, "panel/storage/logs/app.log"),
		MaxSize:  10,
		MaxAge:   30,
		Compress: true,
	}

	level := slog.LevelInfo
	if app.Conf.Bool("app.debug") {
		level = slog.LevelDebug
	}

	app.Logger = slog.New(slog.NewJSONHandler(ljLogger, &slog.HandlerOptions{
		Level: level,
	}))
	slog.SetDefault(app.Logger)
}
