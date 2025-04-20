package logger

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLogLevelString tests the String method of LogLevel
func TestLogLevelString(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{DebugLevel, LogLevelDebug},
		{InfoLevel, LogLevelInfo},
		{WarnLevel, LogLevelWarn},
		{ErrorLevel, LogLevelError},
		{FatalLevel, LogLevelFatal},
		{LogLevel(999), LogLevelUnknown}, // Unknown level
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Level_%d", tt.level), func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.level.String())
		})
	}
}

// TestLogLevelToZapLevel tests the ToZapLevel method of LogLevel
func TestLogLevelToZapLevel(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{DebugLevel, "debug"},
		{InfoLevel, "info"},
		{WarnLevel, "warn"},
		{ErrorLevel, "error"},
		{FatalLevel, "fatal"},
		{LogLevel(999), "info"}, // Default to info for unknown
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Level_%d", tt.level), func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.level.ToZapLevel().String())
		})
	}
}

// CustomTestDestination is a test destination that captures log entries
type CustomTestDestination struct {
	Entries []LogEntry
}

func NewTestDestination() *CustomTestDestination {
	return &CustomTestDestination{
		Entries: make([]LogEntry, 0),
	}
}

func (d *CustomTestDestination) Write(entry LogEntry) error {
	d.Entries = append(d.Entries, entry)
	return nil
}

func (d *CustomTestDestination) Close() error {
	return nil
}

// TestLoggerCreation tests the creation of a logger with options
func TestLoggerCreation(t *testing.T) {
	// Test default logger
	logger := New()
	logger.AddDestination("custom", NewTestDestination())
	assert.Equal(t, "", logger.serviceName)
	assert.Equal(t, InfoLevel, logger.minLevel)
	assert.False(t, logger.isProd)
	assert.Len(t, logger.destinations, 1)
	assert.Contains(t, logger.destinations, "custom")
	assert.Equal(t, []string{}, logger.defaultDests)

	// Test with options
	logger = New(
		WithServiceName("test-service"),
		WithMinLevel(DebugLevel),
		WithProduction(true),
		WithDefaultDestinations("test"),
	)
	assert.Equal(t, "test-service", logger.serviceName)
	assert.Equal(t, InfoLevel, logger.minLevel) // Should be InfoLevel because isProd=true
	assert.True(t, logger.isProd)
	assert.Equal(t, []string{"test"}, logger.defaultDests)

	// Clean up
	logger.Close()
}

// TestLoggerLevelFiltering tests that logs below minimum level are filtered
func TestLoggerLevelFiltering(t *testing.T) {
	testDest := NewTestDestination()

	// Create logger with InfoLevel minimum
	logger := New(
		WithServiceName("test-service"),
		WithMinLevel(InfoLevel),
	)
	logger.AddDestination("test", testDest)
	logger.SetDefaultDestinations("test")

	// Log messages at different levels
	logger.Debug("Debug message", nil)
	logger.Info("Info message", nil)
	logger.Warn("Warn message", nil)

	// Verify only Info and Warn were logged
	assert.Len(t, testDest.Entries, 2)
	assert.Equal(t, "Info message", testDest.Entries[0].Message)
	assert.Equal(t, "Warn message", testDest.Entries[1].Message)

	// Clean up
	logger.Close()
}

// TestLoggerProductionMode tests that debug logs are filtered in production mode
func TestLoggerProductionMode(t *testing.T) {
	testDest := NewTestDestination()

	// Create logger in production mode with DebugLevel
	logger := New(
		WithServiceName("test-service"),
		WithMinLevel(DebugLevel),
		WithProduction(true), // This should override the min level for debug
	)
	logger.AddDestination("test", testDest)
	logger.SetDefaultDestinations("test")

	// Log messages at different levels
	logger.Debug("Debug message", nil)
	logger.Info("Info message", nil)

	// Verify only Info was logged (Debug filtered in production)
	assert.Len(t, testDest.Entries, 1)
	assert.Equal(t, "Info message", testDest.Entries[0].Message)

	// Clean up
	logger.Close()
}

// TestLoggerFields tests that fields are properly included in log entries
func TestLoggerFields(t *testing.T) {
	testDest := NewTestDestination()

	// Create logger
	logger := New(
		WithServiceName("test-service"),
	)
	logger.AddDestination("test", testDest)
	logger.SetDefaultDestinations("test")

	// Log with fields
	fields := map[string]interface{}{
		"string": "value",
		"number": 42,
		"bool":   true,
	}
	logger.Info("Test message", fields)

	// Verify fields are included
	assert.Len(t, testDest.Entries, 1)
	assert.Equal(t, "Test message", testDest.Entries[0].Message)
	assert.Equal(t, fields, testDest.Entries[0].Fields)
	assert.Equal(t, "test-service", testDest.Entries[0].ServiceName)

	// Clean up
	logger.Close()
}

// TestLoggerMultipleDestinations tests logging to multiple destinations
func TestLoggerMultipleDestinations(t *testing.T) {
	testDest1 := NewTestDestination()
	testDest2 := NewTestDestination()

	// Create logger with multiple destinations
	logger := New(
		WithServiceName("test-service"),
	)
	logger.AddDestination("test1", testDest1)
	logger.AddDestination("test2", testDest2)
	logger.SetDefaultDestinations("test1", "test2")

	// Log a message
	logger.Info("Test message", nil)

	// Verify it went to both destinations
	assert.Len(t, testDest1.Entries, 1)
	assert.Len(t, testDest2.Entries, 1)
	assert.Equal(t, "Test message", testDest1.Entries[0].Message)
	assert.Equal(t, "Test message", testDest2.Entries[0].Message)

	// Test explicit destination
	logger.Error("Error message", nil, "test1")

	// Verify it only went to test1
	assert.Len(t, testDest1.Entries, 2)
	assert.Len(t, testDest2.Entries, 1)
	assert.Equal(t, "Error message", testDest1.Entries[1].Message)

	// Clean up
	logger.Close()
}

// TestLoggerFileDestination tests the file destination option
func TestLoggerFileDestination(t *testing.T) {
	// Create a temporary directory for test logs
	tempDir, err := os.MkdirTemp("", "logger_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "test.log")

	// Create logger with file destination
	logger := New(
		WithServiceName("test-service"),
		WithFileDestination(logFile, 10, 1, 1, false),
		WithDefaultDestinations(FileLogger),
	)

	// Log some messages
	logger.Info("Info to file", map[string]interface{}{"key": "value"})
	logger.Errorf("Error to file")

	// Close to ensure everything is flushed
	logger.Close()

	// Verify log file was created and contains the messages
	f, err := os.Open(logFile)
	require.NoError(t, err)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	require.NoError(t, scanner.Err())

	// Should have two log entries
	assert.Len(t, lines, 2)

	// Parse and verify JSON entries
	var infoEntry, errorEntry map[string]interface{}

	err = json.Unmarshal([]byte(lines[0]), &infoEntry)
	require.NoError(t, err)
	assert.Equal(t, "Info to file", infoEntry["msg"])
	assert.Equal(t, "test-service", infoEntry["service"])
	assert.Equal(t, "value", infoEntry["key"])

	err = json.Unmarshal([]byte(lines[1]), &errorEntry)
	require.NoError(t, err)
	assert.Equal(t, "Error to file", errorEntry["msg"])
	assert.Equal(t, "test-service", errorEntry["service"])
}

// TestLoggerAddRemoveDestination tests adding and removing destinations
func TestLoggerAddRemoveDestination(t *testing.T) {
	testDest := NewTestDestination()

	// Create logger
	logger := New()
	logger.AddDestination("test", testDest)
	logger.SetDefaultDestinations("test")

	// Log a message
	logger.Info("Test message", nil)
	assert.Len(t, testDest.Entries, 1)

	// Remove the destination
	logger.RemoveDestination("test")
	logger.Info("Another message", nil)

	// The message shouldn't have been logged to the removed destination
	assert.Len(t, testDest.Entries, 1)

	// Clean up
	logger.Close()
}

// TestConsoleDestination tests the console destination implementation
func TestConsoleDestination(t *testing.T) {
	// Redirect stdout temporarily
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Create console destination
	consoleDest := NewConsoleDestination()

	// Log a test entry
	entry := LogEntry{
		ServiceName: "test-service",
		Level:       InfoLevel,
		Message:     "Console test message",
		Fields:      map[string]interface{}{"test": true},
	}
	err := consoleDest.Write(entry)
	require.NoError(t, err)

	// Close the writer to capture output
	w.Close()

	// Read the output before closing the destination
	var buf strings.Builder
	_, err = io.Copy(&buf, r)
	require.NoError(t, err)

	// Close the destination BEFORE restoring stdout
	// This ensures we sync while the redirect is still active
	_ = consoleDest.Close()

	// Now restore stdout
	os.Stdout = oldStdout

	// Verify output contains the message
	output := buf.String()
	assert.Contains(t, output, "Console test message")
	assert.Contains(t, output, "test-service")
	assert.Contains(t, output, `"test": true`)
}

// TestFatalExit tests that Fatal logs cause program exit
// This test is skipped because it would terminate the test process
func TestFatalExit(t *testing.T) {
	if os.Getenv("TEST_FATAL_EXIT") == "1" {
		logger := New()
		logger.Fatal("This should exit", nil)
		// Should not reach here
		t.Fail()
	} else {
		t.Skip("Skipping fatal exit test as it would terminate the process")
	}
}

// TestWithFileDestinationOption tests the WithFileDestination option
func TestWithFileDestinationOption(t *testing.T) {
	// Create a temporary directory for test logs
	tempDir, err := os.MkdirTemp("", "logger_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "test.log")

	// Create logger with file destination option
	logger := New(
		WithFileDestination(logFile, 10, 1, 1, false),
		WithDefaultDestinations(FileLogger),
	)

	// Verify the file destination was added
	assert.Contains(t, logger.destinations, "file")
	assert.Contains(t, logger.defaultDests, "file")

	// Log a test message
	logger.Info("Test with file destination option", nil)

	// Close to ensure everything is flushed
	logger.Close()

	// Verify log file was created
	_, err = os.Stat(logFile)
	assert.NoError(t, err)
}
