package logger

import (
	"context"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// Logger interface
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
	WithContext(ctx context.Context) Logger
}

// StructuredLogger implements Logger interface
type StructuredLogger struct {
	logger *logrus.Logger
	entry  *logrus.Entry
}

// New creates a new structured logger
func New(level string) Logger {
	log := logrus.New()
	
	// Set log level
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		lvl = logrus.InfoLevel
	}
	log.SetLevel(lvl)
	
	// Set formatter
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
	})
	
	// Set output
	log.SetOutput(os.Stdout)
	
	return &StructuredLogger{
		logger: log,
		entry:  log.WithFields(logrus.Fields{}),
	}
}

func (l *StructuredLogger) Debug(args ...interface{}) {
	l.entry.Debug(args...)
}

func (l *StructuredLogger) Info(args ...interface{}) {
	l.entry.Info(args...)
}

func (l *StructuredLogger) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}

func (l *StructuredLogger) Error(args ...interface{}) {
	l.entry.Error(args...)
}

func (l *StructuredLogger) Fatal(args ...interface{}) {
	l.entry.Fatal(args...)
}

func (l *StructuredLogger) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

func (l *StructuredLogger) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

func (l *StructuredLogger) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

func (l *StructuredLogger) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

func (l *StructuredLogger) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

func (l *StructuredLogger) WithField(key string, value interface{}) Logger {
	return &StructuredLogger{
		logger: l.logger,
		entry:  l.entry.WithField(key, value),
	}
}

func (l *StructuredLogger) WithFields(fields map[string]interface{}) Logger {
	return &StructuredLogger{
		logger: l.logger,
		entry:  l.entry.WithFields(fields),
	}
}

func (l *StructuredLogger) WithContext(ctx context.Context) Logger {
	return &StructuredLogger{
		logger: l.logger,
		entry:  l.entry.WithContext(ctx),
	}
}

// Global logger instance
var globalLogger Logger

// Init initializes global logger
func Init(level string) {
	globalLogger = New(level)
}

// Global logging functions
func Debug(args ...interface{}) {
	if globalLogger == nil {
		fmt.Println("DEBUG:", fmt.Sprint(args...))
		return
	}
	globalLogger.Debug(args...)
}

func Info(args ...interface{}) {
	if globalLogger == nil {
		fmt.Println("INFO:", fmt.Sprint(args...))
		return
	}
	globalLogger.Info(args...)
}

func Warn(args ...interface{}) {
	if globalLogger == nil {
		fmt.Println("WARN:", fmt.Sprint(args...))
		return
	}
	globalLogger.Warn(args...)
}

func Error(args ...interface{}) {
	if globalLogger == nil {
		fmt.Println("ERROR:", fmt.Sprint(args...))
		return
	}
	globalLogger.Error(args...)
}

func Fatal(args ...interface{}) {
	if globalLogger == nil {
		fmt.Println("FATAL:", fmt.Sprint(args...))
		os.Exit(1)
	}
	globalLogger.Fatal(args...)
}
