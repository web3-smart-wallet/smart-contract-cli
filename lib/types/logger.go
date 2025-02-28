package types

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Logger represents a custom logger that writes to both file and stdout
type Logger struct {
	file   *os.File
	logger *log.Logger
}

// NewLogger creates a new logger that writes to file only
func NewLogger() (*Logger, error) {
	// Create logs directory if it doesn't exist
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %v", err)
	}
	// Create log file with timestamp in name
	timestamp := time.Now().Format("2006-01-02")
	logFile := filepath.Join(logsDir, fmt.Sprintf("app_%s.log", timestamp))

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	// Create logger that only writes to file
	logger := log.New(file, "", log.Ldate|log.Ltime)

	return &Logger{
		file:   file,
		logger: logger,
	}, nil
}

// Log writes a message to the log
func (l *Logger) Log(level, message string) {
	l.logger.Printf("[%s] %s", level, message)
}

// Close closes the log file
func (l *Logger) Close() error {
	return l.file.Close()
}
