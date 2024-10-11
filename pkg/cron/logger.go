package cron

import "go.uber.org/zap"

type Logger struct {
	log   *zap.Logger
	debug bool
}

func NewLogger(log *zap.Logger, debug bool) *Logger {
	return &Logger{
		debug: debug,
		log:   log,
	}
}

func (log *Logger) Info(msg string, keysAndValues ...any) {
	if !log.debug {
		return
	}

	log.log.Info(msg, log.toZapFields(keysAndValues...)...)
}

func (log *Logger) Error(err error, msg string, keysAndValues ...any) {
	fields := log.toZapFields(keysAndValues...)
	fields = append(fields, zap.Error(err))
	log.log.Error(msg, fields...)
}

func (log *Logger) toZapFields(keysAndValues ...any) []zap.Field {
	fields := make([]zap.Field, 0, len(keysAndValues)/2)
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			key, okKey := keysAndValues[i].(string)
			value := keysAndValues[i+1]
			if okKey {
				fields = append(fields, zap.Any(key, value))
			}
		}
	}
	return fields
}
