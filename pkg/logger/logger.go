package logger

import "go.uber.org/zap"

type Logger struct {
	client *zap.Logger
}

// Info logs an informational message.
func (logger *Logger) Info(message string) {
	logger.client.Info(message)
}

// Infof logs a formatted informational message.
func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.client.Sugar().Infof(format, args...)
}

// Warn logs a warning message.
func (logger *Logger) Warn(message string) {
	logger.client.Warn(message)
}

// Warnf logs a formatted warning message.
func (logger *Logger) Warnf(format string, args ...interface{}) {
	logger.client.Sugar().Warnf(format, args...)
}

// Error logs an error message.
func (logger *Logger) Error(message string) {
	logger.client.Error(message)
}

// Errorf logs a formatted error message.
func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.client.Sugar().Errorf(format, args...)
}

// Fatal logs a fatal error message and exits the application.
func (logger *Logger) Fatal(message string) {
	logger.client.Fatal(message)
}

// Fatalf logs a formatted fatal error message and exits the application.
func (logger *Logger) Fatalf(format string, args ...interface{}) {
	logger.client.Sugar().Fatalf(format, args...)
}

func NewLogger(component string) Logger {
	logger := zap.Must(zap.NewProduction()).With(
		zap.String("component", component),
	)
	return Logger{
		client: logger,
	}
}
