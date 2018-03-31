package log

import (
	"github.com/deepak11627/arc/arc"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config contains configuration options for creating the logger
type Config struct {
	// JSONFormat determines if logs are logged in JSON
	JSONFormat bool
	// Debug determines if debug level logs are logged. Logging levels may be modified at runtime via a Manager
	Debug bool
	// LogPath specifies the path to the log file. stdout and stderr are also acceptable values
	LogPath string
	// ApplicationVersion specifies the version of the application, to be imprinted on all log messages
	ApplicationVersion string
}

const (
	consoleEncoding  = "console"
	jsonEncoding     = "json"
	stdoutOutputPath = "stdout"

	timeKey  = "time"
	levelKey = "loglevel"
)

// Logger used for loggin purpose through out the application
type Logger struct {
	*zap.SugaredLogger
}

// NewLogger return a logger
func NewLogger(config *Config) (arc.Logger, error) {
	level := zap.InfoLevel
	if config.Debug {
		level = zap.DebugLevel
	}

	encoding := consoleEncoding
	if config.JSONFormat {
		encoding = jsonEncoding
	}

	outputPath := stdoutOutputPath
	if config.LogPath != "" {
		outputPath = config.LogPath
	}

	cfg := zap.NewProductionConfig()

	cfg.Level = zap.NewAtomicLevelAt(level)
	cfg.Encoding = encoding
	cfg.OutputPaths = []string{outputPath}

	cfg.EncoderConfig.TimeKey = timeKey
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.LevelKey = levelKey
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	logger, err := cfg.Build(zap.Fields(
		zap.String("service", "ARC-IN-GO"),
		zap.String("version", config.ApplicationVersion),
	))

	if err != nil {
		return nil, err
	}

	return &Logger{
		logger.Sugar(),
	}, nil

}

// Debug logs a debug message to the zap logger
func (l *Logger) Debug(msg string, keyvals ...interface{}) {
	l.Debugw(msg, keyvals...)
}

// Info logs an info message to the zap logger
func (l *Logger) Info(msg string, keyvals ...interface{}) {
	l.Infow(msg, keyvals...)
}

// Warn logs a warning message to the zap logger
func (l *Logger) Warn(msg string, keyvals ...interface{}) {
	l.Warnw(msg, keyvals...)
}

// Error logs a error message to the zap logger
func (l *Logger) Error(msg string, keyvals ...interface{}) {
	l.Errorw(msg, keyvals...)
}
