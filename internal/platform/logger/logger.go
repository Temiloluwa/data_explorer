package logger

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger interface defines the logging methods used in the application.
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	// With returns a new logger with the specified structured context added.
	With(args ...interface{}) Logger
}

// Config holds logging configuration.
type Config struct {
	Level  string `mapstructure:"level"`  // e.g., "debug", "info", "warn", "error", "fatal"
	Format string `mapstructure:"format"` // e.g., "console", "json" (defaults to console)
	// Add more options like output path if needed
}

// zapLogger is an implementation of the Logger interface using Uber's Zap.
type zapLogger struct {
	sugar *zap.SugaredLogger
}

// New creates a new Logger instance based on the provided configuration.
func New(config Config) (Logger, error) {
	level := zap.NewAtomicLevel()
	if err := level.UnmarshalText([]byte(strings.ToLower(config.Level))); err != nil {
		// Default to Info level if parsing fails
		fmt.Fprintf(os.Stderr, "Warning: Invalid log level '%s', defaulting to 'info'. Error: %v\n", config.Level, err)
		level.SetLevel(zapcore.InfoLevel)
	}

	// Configure encoder based on format
	var encoder zapcore.Encoder
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp" // Use "ts" for shorter keys if preferred
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder // e.g., DEBUG, INFO
	encoderCfg.MessageKey = "message"
	encoderCfg.CallerKey = "caller" // Display caller file/line

	if strings.ToLower(config.Format) == "json" {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		// Default to console format
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder // Add colors for console
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	// Configure core (output destination and level)
	// TODO: Add option to write to a file based on config.OutputPath
	writer := zapcore.Lock(os.Stdout) // Write to stdout
	core := zapcore.NewCore(encoder, writer, level)

	// Create the base Zap logger
	// AddCallerSkip(1) adjusts the caller info to show the actual caller of Debugf/Infof etc.
	// instead of the zapLogger methods themselves.
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel)) // Add stacktrace for Error level and above

	// Use SugaredLogger for convenience (printf-style formatting)
	sugar := logger.Sugar()

	// Optionally sync logger on exit (recommended)
	// defer sugar.Sync() // This should be called at application exit, not here.

	return &zapLogger{sugar: sugar}, nil
}

func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.sugar.Debugf(format, args...)
}

func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.sugar.Infof(format, args...)
}

func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.sugar.Warnf(format, args...)
}

func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.sugar.Errorf(format, args...)
}

func (l *zapLogger) Fatalf(format string, args ...interface{}) {
	l.sugar.Fatalf(format, args...)
}

// With creates a new logger instance with the added context fields.
func (l *zapLogger) With(args ...interface{}) Logger {
	newSugar := l.sugar.With(args...)
	return &zapLogger{sugar: newSugar}
}

// Ensure zapLogger implements Logger at compile time.
var _ Logger = (*zapLogger)(nil)
