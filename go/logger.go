package main

import (
	"context"
	"log/slog"
	"os"
)

// Logger wraps the structured logger with application-specific methods
type Logger struct {
	*slog.Logger
}

// NewLogger creates a new structured logger
func NewLogger() *Logger {
	// Use JSON handler for structured logging
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	
	return &Logger{
		Logger: slog.New(handler),
	}
}

// WithContext adds context fields to the logger
func (l *Logger) WithContext(ctx context.Context) *slog.Logger {
	// Add any context-specific fields here if needed
	return l.Logger
}

// WithRequest adds request-specific fields to the logger
func (l *Logger) WithRequest(method, path string) *slog.Logger {
	return l.Logger.With(
		"method", method,
		"path", path,
	)
}

// WithError adds error information to the logger
func (l *Logger) WithError(err error) *slog.Logger {
	return l.Logger.With("error", err.Error())
}

// WithFields adds arbitrary fields to the logger
func (l *Logger) WithFields(fields map[string]any) *slog.Logger {
	args := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return l.Logger.With(args...)
}