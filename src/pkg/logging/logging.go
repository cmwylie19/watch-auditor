package logging

import (
	"log/slog"
	"os"
)

// LoggerInterface defines the methods that the logger should implement.
// This allows for flexibility in using different types of loggers, such as mocks in tests.
type LoggerInterface interface {
	Info(msg string)
	Debug(msg string)
	Warn(msg string)
	Error(msg string)
	SetLevel(level slog.Level)
	CloseFile()
}

// Ensure that Logger implements LoggerInterface.
var _ LoggerInterface = (*Logger)(nil)

// Logger wraps a slog.Logger to provide structured logging.
type Logger struct {
	logger   *slog.Logger
	logLevel *slog.LevelVar
	file     *os.File
}

// NewLogger creates a new Logger instance with the specified output file.
// If filePath is empty, logs will be written to stdout.
func NewLogger(filePath string) (*Logger, error) {
	logLevel := &slog.LevelVar{}
	logLevel.Set(slog.LevelInfo) // Default level is INFO

	var handler slog.Handler
	var file *os.File
	var err error

	if filePath != "" {
		file, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		handler = slog.NewJSONHandler(file, &slog.HandlerOptions{Level: logLevel})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	}

	l := slog.New(handler)

	return &Logger{
		logger:   l,
		logLevel: logLevel,
		file:     file,
	}, nil
}

// CloseFile closes the file if it was opened.
func (l *Logger) CloseFile() {
	if l.file != nil {
		l.file.Close()
	}
}

// SetLevel adjusts the logging level.
func (l *Logger) SetLevel(level slog.Level) {
	l.logLevel.Set(level)
}

// Info logs an info message.
func (l *Logger) Info(msg string) {
	l.logger.Info(msg)
}

// Debug logs a debug message.
func (l *Logger) Debug(msg string) {
	l.logger.Debug(msg)
}

// Warn logs a warning message.
func (l *Logger) Warn(msg string) {
	l.logger.Warn(msg)
}

// Error logs an error message.
func (l *Logger) Error(msg string) {
	l.logger.Error(msg)
}

// SetCustomLogger allows setting a custom logger, useful for testing.
func (l *Logger) SetCustomLogger(customLogger *slog.Logger) {
	l.logger = customLogger
}
