package scmlogger

var logger Logger

type Logger interface {
	Log(msg string, args ...interface{})
}

func GetLogger() Logger {
	if logger == nil {
		logger = &noopLogger{}
	}
	return logger
}

type noopLogger struct{}

func (l *noopLogger) Log(_ string, _ ...interface{}) {}
