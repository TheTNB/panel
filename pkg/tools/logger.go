package tools

// NoopLogger 给 go-resty 使用的空日志
type NoopLogger struct{}

func (NoopLogger) Errorf(format string, v ...any) {}

func (NoopLogger) Warnf(format string, v ...any) {}

func (NoopLogger) Debugf(format string, v ...any) {}
