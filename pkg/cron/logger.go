package cron

import (
	"log/slog"
)

type Logger struct {
	log   *slog.Logger
	debug bool
}

func NewLogger(log *slog.Logger, debug bool) *Logger {
	return &Logger{
		debug: debug,
		log:   log,
	}
}

func (log *Logger) Info(msg string, keysAndValues ...any) {
	if !log.debug {
		return
	}

	log.log.Info(msg, keysAndValues...)
}

func (log *Logger) Error(err error, msg string, keysAndValues ...any) {
	fields := []any{slog.Any("err", err)}
	fields = append(fields, log.toSlogArgs(keysAndValues...)...)
	log.log.Error(msg, fields...)
}

func (log *Logger) toSlogArgs(keysAndValues ...any) []any {
	fields := make([]any, 0, len(keysAndValues)/2)
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			key, ok := keysAndValues[i].(string)
			if ok {
				fields = append(fields, slog.Any(key, keysAndValues[i+1]))
			}
		}
	}
	return fields
}
