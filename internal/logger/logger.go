package logger

import (
	"log"
	"os"
	"sync"
)

var (
	logger *Logger
	once   sync.Once
)

// Logger represents a logger instance
type Logger struct {
	debug *log.Logger
	info  *log.Logger
	error *log.Logger
	isDebug bool
}

// Init initializes the logger
func Init(debug bool) {
	once.Do(func() {
		logger = &Logger{
			debug: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
			info:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
			error: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
			isDebug: debug,
		}
	})
}

// GetLogger returns the logger instance
func GetLogger() *Logger {
	if logger == nil {
		Init(false)
	}
	return logger
}

// Debug logs a debug message
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.isDebug {
		l.debug.Printf(format, v...)
	}
}

// Info logs an info message
func (l *Logger) Info(format string, v ...interface{}) {
	l.info.Printf(format, v...)
}

// Error logs an error message
func (l *Logger) Error(format string, v ...interface{}) {
	l.error.Printf(format, v...)
}

// Fatal logs an error message and exits
func (l *Logger) Fatal(format string, v ...interface{}) {
	l.error.Fatalf(format, v...)
}