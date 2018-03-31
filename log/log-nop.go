package log

import (
	"github.com/deepak11627/arc/arc"
	"go.uber.org/zap"
)

// NopLogger for testing purpose
type NopLogger struct {
	*zap.Logger
}

// Debug logs a debug message to the zap logger
func (l *NopLogger) Debug(msg string, keyvals ...interface{}) {
	//	l.Debug(msg, keyvals...)
	return
}

// Info logs an info message to the zap NopLogger
func (l *NopLogger) Info(msg string, keyvals ...interface{}) {
	//l.Info(msg, keyvals...)
	return
}

// Warn logs a warning message to the zap NopLogger
func (l *NopLogger) Warn(msg string, keyvals ...interface{}) {
	//l.Warn(msg, keyvals...)
	return
}

// Error logs a error message to the zap NopLogger
func (l *NopLogger) Error(msg string, keyvals ...interface{}) {
	//l.Error(msg, keyvals...)
	return
}

// NewNop logger
func NewNopLogger() arc.Logger {
	return &NopLogger{
		zap.NewNop(),
	}
}
