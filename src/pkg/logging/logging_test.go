package logging

import (
	"bytes"
	"log/slog"
	"os"
	"testing"
)

// TestLogger_NewLogger_NoFilePath tests the logger when no file path is provided.
func TestLogger_NewLogger_NoFilePath(t *testing.T) {
	logger, err := NewLogger("")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Set log level to INFO and test logging to stdout
	logger.SetLevel(slog.LevelInfo)
	logger.Info("This is an info message")
}

// TestLogger_NewLogger_FileError tests the logger initialization when a file error occurs.
func TestLogger_NewLogger_FileError(t *testing.T) {
	// Provide an invalid file path to simulate an error
	_, err := NewLogger("/invalid_path/test_log.json")
	if err == nil {
		t.Fatalf("Expected an error due to invalid file path, but got none")
	}
}

// TestLogger_SetLevel tests the setting of various log levels.
func TestLogger_SetLevel(t *testing.T) {
	filePath := "test_log.json"
	defer os.Remove(filePath)

	logger, err := NewLogger(filePath)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.CloseFile()

	// Test setting DEBUG level
	logger.SetLevel(slog.LevelDebug)
	logger.Debug("This is a debug message")
	validateLogFile(t, filePath, `"level":"DEBUG"`)
	clearFile(filePath)

	// Test setting INFO level
	logger.SetLevel(slog.LevelInfo)
	logger.Info("This is an info message")
	validateLogFile(t, filePath, `"level":"INFO"`)
	clearFile(filePath)

	// Test setting WARN level
	logger.SetLevel(slog.LevelWarn)
	logger.Warn("This is a warning message")
	validateLogFile(t, filePath, `"level":"WARN"`)
	clearFile(filePath)

	// Test setting ERROR level
	logger.SetLevel(slog.LevelError)
	logger.Error("This is an error message")
	validateLogFile(t, filePath, `"level":"ERROR"`)
}

// TestLogger_Info tests the Info logging method.
func TestLogger_Info(t *testing.T) {
	filePath := "test_log.json"
	defer os.Remove(filePath)

	logger, err := NewLogger(filePath)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.CloseFile()

	logger.Info("Info log message")
	validateLogFile(t, filePath, `"level":"INFO"`)
	validateLogFile(t, filePath, `"msg":"Info log message"`)
}

// TestLogger_Debug tests the Debug logging method.
func TestLogger_Debug(t *testing.T) {
	filePath := "test_log.json"
	defer os.Remove(filePath)

	logger, err := NewLogger(filePath)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.CloseFile()

	logger.SetLevel(slog.LevelDebug)
	logger.Debug("Debug log message")
	validateLogFile(t, filePath, `"level":"DEBUG"`)
	validateLogFile(t, filePath, `"msg":"Debug log message"`)
}

// TestLogger_Warn tests the Warn logging method.
func TestLogger_Warn(t *testing.T) {
	filePath := "test_log.json"
	defer os.Remove(filePath)

	logger, err := NewLogger(filePath)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.CloseFile()

	logger.Warn("Warn log message")
	validateLogFile(t, filePath, `"level":"WARN"`)
	validateLogFile(t, filePath, `"msg":"Warn log message"`)
}

// TestLogger_Error tests the Error logging method.
func TestLogger_Error(t *testing.T) {
	filePath := "test_log.json"
	defer os.Remove(filePath)

	logger, err := NewLogger(filePath)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.CloseFile()

	logger.Error("Error log message")
	validateLogFile(t, filePath, `"level":"ERROR"`)
	validateLogFile(t, filePath, `"msg":"Error log message"`)
}

// TestLogger_SetCustomLogger tests setting a custom logger.
func TestLogger_SetCustomLogger(t *testing.T) {
	filePath := "test_log.json"
	defer os.Remove(filePath)

	logger, err := NewLogger(filePath)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.CloseFile()

	// Create a custom logger that logs to a buffer
	var buf bytes.Buffer
	customLogger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{}))
	logger.SetCustomLogger(customLogger)

	// Log a message with the custom logger
	logger.Info("Custom log message")

	// Verify that the message was logged to the buffer
	if !contains(buf.Bytes(), `"msg":"Custom log message"`) {
		t.Errorf("Expected buffer to contain 'Custom log message', but got %s", buf.String())
	}
}

// validateLogFile checks if the log file contains the expected content.
func validateLogFile(t *testing.T, filePath, expected string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	if !contains(content, expected) {
		t.Errorf("Expected log file to contain %s, but got %s", expected, string(content))
	}
}

// contains checks if content contains the expected substring.
func contains(content []byte, expected string) bool {
	return bytes.Contains(content, []byte(expected))
}

// clearFile clears the contents of the file.
func clearFile(filePath string) {
	os.Truncate(filePath, 0)
}
