package monitoring

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap logger with service context
type Logger struct {
	*zap.Logger
	serviceName string
}

// LoggerConfig represents logger configuration
type LoggerConfig struct {
	ServiceName string
	Level       string // debug, info, warn, error
	Format      string // json, console
	OutputPath  string // stdout, stderr, or file path
}

// NewLogger creates a new structured logger
func NewLogger(config LoggerConfig) (*Logger, error) {
	var zapConfig zap.Config

	switch config.Format {
	case "console":
		zapConfig = zap.NewDevelopmentConfig()
	default:
		zapConfig = zap.NewProductionConfig()
	}

	// Set log level
	level, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}
	zapConfig.Level = zap.NewAtomicLevelAt(level)

	// Set output path
	if config.OutputPath != "" {
		zapConfig.OutputPaths = []string{config.OutputPath}
	}

	logger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	// Add service name to all logs
	logger = logger.With(zap.String("service", config.ServiceName))

	return &Logger{
		Logger:      logger,
		serviceName: config.ServiceName,
	}, nil
}

// WithRequest adds request context to logger
func (l *Logger) WithRequest(requestID string) *zap.Logger {
	return l.With(zap.String("request_id", requestID))
}

// WithError adds error context to logger
func (l *Logger) WithError(err error) *zap.Logger {
	return l.With(zap.Error(err))
}

// HTTP request logging
func (l *Logger) LogHTTPRequest(method, path string, statusCode int, duration int64) {
	l.Info("HTTP request",
		zap.String("method", method),
		zap.String("path", path),
		zap.Int("status_code", statusCode),
		zap.Int64("duration_ms", duration),
	)
}

// Database operation logging
func (l *Logger) LogDBOperation(operation, table string, duration int64, err error) {
	if err != nil {
		l.Error("Database operation failed",
			zap.String("operation", operation),
			zap.String("table", table),
			zap.Int64("duration_ms", duration),
			zap.Error(err),
		)
	} else {
		l.Debug("Database operation",
			zap.String("operation", operation),
			zap.String("table", table),
			zap.Int64("duration_ms", duration),
		)
	}
}

// External service call logging
func (l *Logger) LogExternalCall(service, endpoint string, statusCode int, duration int64, err error) {
	if err != nil {
		l.Error("External service call failed",
			zap.String("external_service", service),
			zap.String("endpoint", endpoint),
			zap.Int("status_code", statusCode),
			zap.Int64("duration_ms", duration),
			zap.Error(err),
		)
	} else {
		l.Info("External service call",
			zap.String("external_service", service),
			zap.String("endpoint", endpoint),
			zap.Int("status_code", statusCode),
			zap.Int64("duration_ms", duration),
		)
	}
}
