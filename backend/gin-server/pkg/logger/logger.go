// package logger provides a configurable logging solution
package logger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Log levels as strings
const (
	LogLevelDebug   = "DEBUG"
	LogLevelInfo    = "INFO"
	LogLevelWarn    = "WARN"
	LogLevelError   = "ERROR"
	LogLevelFatal   = "FATAL"
	LogLevelUnknown = "UNKNOWN"
)

// LogLevel represents the severity level of a log message
type LogLevel int

const (
	// DebugLevel logs detailed information for debugging
	DebugLevel LogLevel = iota
	// InfoLevel logs general operational information
	InfoLevel
	// WarnLevel logs potentially harmful situations
	WarnLevel
	// ErrorLevel logs error events that might still allow the application to continue
	ErrorLevel
	// FatalLevel logs severe error events that will lead the application to abort
	FatalLevel
)

const (
	FileLogger    = "file"
	ConsoleLogger = "console"
)

// String returns string representation of log level
func (l LogLevel) String() string {
	switch l {
	case DebugLevel:
		return LogLevelDebug
	case InfoLevel:
		return LogLevelInfo
	case WarnLevel:
		return LogLevelWarn
	case ErrorLevel:
		return LogLevelError
	case FatalLevel:
		return LogLevelFatal
	default:
		return LogLevelUnknown
	}
}

// ToZapLevel converts our LogLevel to zapcore.Level
func (l LogLevel) ToZapLevel() zapcore.Level {
	switch l {
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case FatalLevel:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// Destination represents where logs can be written
type Destination interface {
	// Write takes a message and writes it to the destination
	Write(entry LogEntry) error
	// Close closes the destination
	Close() error
}

// LogEntry represents a single log message
type LogEntry struct {
	ServiceName string
	Level       LogLevel
	Message     string
	Fields      map[string]interface{}
}

// ConsoleDestination writes logs to console
type ConsoleDestination struct {
	logger *zap.Logger
}

// NewConsoleDestination creates a new console destination
func NewConsoleDestination() *ConsoleDestination {
	config := zap.NewDevelopmentEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.TimeKey = "timestamp"

	consoleEncoder := zapcore.NewConsoleEncoder(config)
	core := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
	logger := zap.New(core)

	return &ConsoleDestination{
		logger: logger,
	}
}

// Write implements Destination
func (c *ConsoleDestination) Write(entry LogEntry) error {
	fields := make([]zap.Field, 0, len(entry.Fields)+1)
	fields = append(fields, zap.String("service", entry.ServiceName))

	for k, v := range entry.Fields {
		fields = append(fields, zap.Any(k, v))
	}

	switch entry.Level {
	case DebugLevel:
		c.logger.Debug(entry.Message, fields...)
	case InfoLevel:
		c.logger.Info(entry.Message, fields...)
	case WarnLevel:
		c.logger.Warn(entry.Message, fields...)
	case ErrorLevel:
		c.logger.Error(entry.Message, fields...)
	case FatalLevel:
		c.logger.Fatal(entry.Message, fields...)
	}

	return nil
}

// Close implements Destination
func (c *ConsoleDestination) Close() error {
	// Ignore ENOTTY error which occurs when stdout is not a terminal
	// READ-MORE: https://github.com/uber-go/zap/issues/991#issuecomment-962098428
	if err := c.logger.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		return err
	}
	return nil
}

// FileDestination writes logs to a file
type FileDestination struct {
	logger     *zap.Logger
	lumberjack *lumberjack.Logger
}

// NewFileDestination creates a new file destination
func NewFileDestination(path string, maxSize int, maxBackups int, maxAge int, compress bool) *FileDestination {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Failed to create log directory: %v\n", err)
	}

	lumberjackLogger := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    maxSize,    // megabytes
		MaxBackups: maxBackups, // number of backups
		MaxAge:     maxAge,     // days
		Compress:   compress,   // compress backups
	}

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.TimeKey = "timestamp"

	fileEncoder := zapcore.NewJSONEncoder(config)

	core := zapcore.NewCore(fileEncoder, zapcore.AddSync(lumberjackLogger), zapcore.DebugLevel)
	logger := zap.New(core)

	return &FileDestination{
		logger:     logger,
		lumberjack: lumberjackLogger,
	}
}

// Write implements Destination
func (f *FileDestination) Write(entry LogEntry) error {
	fields := make([]zap.Field, 0, len(entry.Fields)+1)
	fields = append(fields, zap.String("service", entry.ServiceName))

	for k, v := range entry.Fields {
		fields = append(fields, zap.Any(k, v))
	}

	switch entry.Level {
	case DebugLevel:
		f.logger.Debug(entry.Message, fields...)
	case InfoLevel:
		f.logger.Info(entry.Message, fields...)
	case WarnLevel:
		f.logger.Warn(entry.Message, fields...)
	case ErrorLevel:
		f.logger.Error(entry.Message, fields...)
	case FatalLevel:
		f.logger.Fatal(entry.Message, fields...)
	}

	return nil
}

// Close implements Destination
func (f *FileDestination) Close() error {
	if err := f.logger.Sync(); err != nil {
		return err
	}
	return f.lumberjack.Close()
}

// Logger is the main logger interface
type Logger struct {
	serviceName  string
	minLevel     LogLevel
	isProd       bool
	destinations map[string]Destination
	defaultDests []string
	mu           sync.RWMutex
}

// Option defines a function signature for configuration options
type Option func(*Logger)

// WithServiceName sets the service name
func WithServiceName(name string) Option {
	return func(l *Logger) {
		l.serviceName = name
	}
}

// WithMinLevel sets the minimum log level
func WithMinLevel(level LogLevel) Option {
	return func(l *Logger) {
		l.minLevel = level
	}
}

// WithProduction sets the logger to production mode
func WithProduction(isProd bool) Option {
	return func(l *Logger) {
		l.isProd = isProd
	}
}

// WithDefaultDestinations sets the default destinations
func WithDefaultDestinations(dests ...string) Option {
	return func(l *Logger) {
		l.defaultDests = dests
	}
}

// WithConsoleDestination adds a console destination
func WithConsoleDestination() Option {
	return func(l *Logger) {
		l.destinations[ConsoleLogger] = NewConsoleDestination()
	}
}

// WithFileDestination adds a file destination
func WithFileDestination(path string, maxSize, maxBackups, maxAge int, compress bool) Option {
	return func(l *Logger) {
		l.destinations[FileLogger] = NewFileDestination(path, maxSize, maxBackups, maxAge, compress)
	}
}

// New creates a new logger with the given options
func New(options ...Option) *Logger {
	l := &Logger{
		serviceName:  "",
		minLevel:     InfoLevel,
		isProd:       false,
		destinations: make(map[string]Destination),
		defaultDests: []string{},
	}

	// Apply options
	for _, option := range options {
		option(l)
	}

	// If in production mode, set minimum level to info
	if l.isProd && l.minLevel < InfoLevel {
		l.minLevel = InfoLevel
	}

	return l
}

// log sends the log message to specified destinations
func (l *Logger) log(level LogLevel, msg string, fields map[string]interface{}, dests ...string) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	// Skip if level is below minimum (especially for debug in prod)
	if level < l.minLevel {
		return
	}

	// If no destinations specified, use defaults
	if len(dests) == 0 {
		dests = l.defaultDests
	}

	entry := LogEntry{
		ServiceName: l.serviceName,
		Level:       level,
		Message:     msg,
		Fields:      fields,
	}

	// Write to all specified destinations
	for _, destName := range dests {
		if dest, ok := l.destinations[destName]; ok {
			// Just log destination write errors to stderr for now
			if err := dest.Write(entry); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to write log to destination %s: %v\n", destName, err)
			}
		}
	}
}

// AddDestination adds a destination to the logger
func (l *Logger) AddDestination(name string, dest Destination) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.destinations[name] = dest
}

// RemoveDestination removes a destination from the logger
func (l *Logger) RemoveDestination(name string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if dest, ok := l.destinations[name]; ok {
		_ = dest.Close()
		delete(l.destinations, name)
	}
}

// SetDefaultDestinations sets the default destinations
func (l *Logger) SetDefaultDestinations(dests ...string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.defaultDests = dests
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields map[string]interface{}, dests ...string) {
	l.log(DebugLevel, msg, fields, dests...)
}

// Info logs an info message
func (l *Logger) Info(msg string, fields map[string]interface{}, dests ...string) {
	l.log(InfoLevel, msg, fields, dests...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields map[string]interface{}, dests ...string) {
	l.log(WarnLevel, msg, fields, dests...)
}

// Error logs an error message
func (l *Logger) Error(msg string, fields map[string]interface{}, dests ...string) {
	l.log(ErrorLevel, msg, fields, dests...)
}

// Fatal logs a fatal message and terminates the program
func (l *Logger) Fatal(msg string, fields map[string]interface{}, dests ...string) {
	l.log(FatalLevel, msg, fields, dests...)
	os.Exit(1)
}

// Debug logs a debug message
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log(DebugLevel, fmt.Sprintf(format, args...), nil)
}

// Info logs an info message
func (l *Logger) Infof(format string, args ...interface{}) {
	l.log(InfoLevel, fmt.Sprintf(format, args...), nil)
}

// Warn logs a warning message
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log(WarnLevel, fmt.Sprintf(format, args...), nil)
}

// Error logs an error message
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log(ErrorLevel, fmt.Sprintf(format, args...), nil)
}

// Fatal logs a fatal message and terminates the program
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.log(FatalLevel, fmt.Sprintf(format, args...), nil)
	os.Exit(1)
}

// Close closes all destinations
func (l *Logger) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()

	for name, dest := range l.destinations {
		if err := dest.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close destination %s: %v\n", name, err)
		}
	}
}
