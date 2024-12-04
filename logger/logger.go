package logger

import (
	"fmt"
	"os"
	"time"
)

// logFilePath is the path to the log file
const logFilePath = "/home/cursework/logs"

// Log writes a message to the log file with a timestamp
func Log(action string, details string) error {
	// Ensure the logs directory exists
	logDir := "/home/cursework"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create log directory: %w", err)
		}
	}

	// Open the log file for appending (create it if it doesn't exist)
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	// Write the log entry
	timestamp := time.Now().Format("2006-01-02 03:04:05")
	logEntry := fmt.Sprintf("[%s] %s: %s\n", timestamp, action, details)
	_, err = file.WriteString(logEntry)
	return err
}
