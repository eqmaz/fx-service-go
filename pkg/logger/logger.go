package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// Logger is a wrapper around logrus.Logger
type Logger struct {
	*logrus.Logger
}

// NewLogger creates a new logger with the default settings
func NewLogger() *Logger {
	log := logrus.New()
	log.Out = os.Stdout

	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})

	return &Logger{log}
}

// WithContextFields returns a new logger with the given fields added to the context
func (l *Logger) WithContextFields(fields map[string]interface{}) *Logger {
	return &Logger{l.WithFields(fields).Logger}
}

// Info logs a message with the Info level
func (l *Logger) Info(msg string, fields map[string]interface{}) {
	l.WithFields(fields).Log(logrus.InfoLevel, msg)
}

// Error logs a message with the Error level
func (l *Logger) Error(msg string, fields map[string]interface{}) {
	l.WithFields(fields).Log(logrus.ErrorLevel, msg)
}

// Warn logs a message with the Warn level
func (l *Logger) Warn(msg string, fields map[string]interface{}) {
	l.WithFields(fields).Log(logrus.WarnLevel, msg)
}

// Debug logs a message with the Debug level
func (l *Logger) Debug(msg string, fields map[string]interface{}) {
	l.WithFields(fields).Log(logrus.DebugLevel, msg)
}
